package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MapWidth  = 10
	MapHeight = 8
)

type GameScore struct {
	Name string
	Gold int
}

type Cell struct {
	X int
	Y int
}

type PlayerInfo struct {
	Name             string
	X                int
	Y                int
	Gold             int
	Rank             float32 // 排名比例(当前排名/总人数), 越低越好
	Joined           bool    // 是否参加了本局比赛
	EstimateLeft     int     // 预测的下一步到达后剩余
	Bingo            int     // 预测准确的回合数
	Precision        float32 // 预测准确率
	HighestPrecision float32 // 历史最高准确率
	LowestPrecision  float32 // 历史最低准确率
	TotalPrecision   float32 // 总预测准确率
	TotalGame        int     // 总参与局数
	AverPrecision    float32 // 平均准确率
	Cells            []*Cell
}

/*
	msgtype 0登陆 -1登陆失败 1请准备 2准备好了 3游戏信息 4玩家行动 5游戏结束
*/
type Msg struct {
	Msgtype int
	Token   string
	RoundID int
	X       int
	Y       int
	Sorted  []*GameScore `json:"Results,omitempty"`
}

type Tile struct {
	Gold             int
	Players          []*GameScore `json:"Players,omitempty"`
	x                int          // 0 ~ MapWidth - 1
	y                int          // 0 ~ MapHeight - 1
	cost             int          // 我移动过来的消耗
	left             int          // 我移动过来的剩余
	leftFlag         int          // left是否计算过了
	reachable        int          // 这个格子上可能的人数(除了我自己)
	opportunity      float32      // 机遇(到达这个位置后金币数比自己到达这个位置后金币数多的人数)
	risk             float32      // 风险(到达这个位置后金币数比自己到达这个位置后金币数少的人数)
	totalLeft        int          // 能到达这个位置的所有玩家剩余金币总和(用来算平均分)
	totalOpportunity int          // 所有到达这个位置后金币比我多的人的金币总和
	totalRisk        int          // 所有到达这个位置后金币比我少的人的金币综合
	expected         float32      // 预期收益
}

type Game struct {
	GameID  uint64
	Msgtype int
	status  int //-1无效0准备1开始
	RoundID int
	Wid     int
	Hei     int

	Tilemap [MapHeight][MapWidth]*Tile

	roundRecords []string
	Sorted       []*GameScore `json:"Results,omitempty"`
}

type Record struct {
	RoundId   int
	X         int
	Y         int
	Gold      int
	Expected  float32
	TargetX   int
	TargetY   int
	Crowded   int
	Rank      float32
	FirstGold int
	FirstName string
}

func saveGame(gameID uint64) error {
	filename := fmt.Sprintf("./log%d-%d-%d-%d.txt", time.Now().YearDay(), time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	gameStr := fmt.Sprintf("Game: %d\n", gameID)
	file.Write([]byte(gameStr))
	for i := 0; i < len(records); i++ {
		bytes, e := json.Marshal(records[i])
		if e == nil {
			file.Write(bytes)
			file.Write([]byte("\n"))
		} else {
			return e
		}
	}

	for p := range allPlayers {
		player, _ := allPlayers[p]
		if player.Joined {
			player.TotalGame++
			player.TotalPrecision += player.Precision
			player.AverPrecision = player.TotalPrecision / float32(player.TotalGame)
			if player.Precision > player.HighestPrecision {
				player.HighestPrecision = player.Precision
			}
			if player.LowestPrecision > player.Precision {
				player.LowestPrecision = player.Precision
			}
			bytes, e := json.Marshal(player)
			if e == nil {
				file.Write(bytes)
				file.Write([]byte("\n"))
			} else {
				return e
			}
		}

		player.Joined = false
	}

	return nil
}

func sendMessage(data interface{}) error {
	if ws != nil {
		return ws.WriteJSON(data)
	}
	return errors.New("ws is nil")
}

func recvMessage(data interface{}) error {
	if ws != nil {
		return ws.ReadJSON(data)
	} else {
		return errors.New("ws == nil")
	}
}

func login(uri, token string) error {
	if ws != nil {
		ws.Close()
	}

	var err error
	ws, _, err = websocket.DefaultDialer.Dial(uri, nil)
	if err != nil {
		fmt.Println("websocket连接出错, err = ", err)
		time.Sleep(time.Second)
		return err
	}

	loginMsg := Msg{}
	loginMsg.Msgtype = 0
	loginMsg.Token = token
	err = sendMessage(loginMsg)
	if err != nil {
		fmt.Println("Send error while login msg, err = ", err)
		time.Sleep(time.Second)
		return err
	}

	for {
		time.Sleep(time.Millisecond * 10)
		data := Msg{}
		err = recvMessage(&data)
		if err == nil {
			if data.Msgtype == -1 {
				time.Sleep(time.Second * 5)
				fmt.Println("服务器拒绝, 登录失败")
				return errors.New("login failed")
			}
			fmt.Println("登录成功, msg = ", data)
			return nil
		} else {
			return err
		}
	}
}

func sendMove(x, y, roundId int) error {
	if x < 0 {
		x = 0
	} else if x >= MapWidth {
		x = MapWidth - 1
	}

	if y < 0 {
		y = 0
	} else if y >= MapHeight {
		y = MapHeight - 1
	}

	msg := Msg{}
	msg.Msgtype = 4
	msg.Token = myToken
	msg.RoundID = roundId
	msg.X = x
	msg.Y = y
	return sendMessage(msg)
}

func updateFrame(frameData Game) {
	// 捕获处理过程中的异常, 保证不会出现闪退
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(time.Now(), "捕获到异常: ", r)
		}
	}()

	roundId := frameData.RoundID
	tileMap := frameData.Tilemap
	// 整理出所有的玩家信息
	players := make([]*PlayerInfo, 0, 8)
	if frameData.RoundID == 0 {
		// 第一局, 需要处理新玩家加入, 还有标记哪些玩家进入了游戏
		for i := 0; i < len(tileMap); i++ {
			for j := 0; j < len(tileMap[i]); j++ {
				tileMap[i][j].x = j
				tileMap[i][j].y = i
				if len(tileMap[i][j].Players) > 0 {
					playerLen := len(tileMap[i][j].Players)
					for k := 0; k < playerLen; k++ {
						p := tileMap[i][j].Players[k]
						player, ok := allPlayers[p.Name]
						if ok {
							player.Gold = p.Gold
							player.X = myX
							player.Y = myY
							player.Joined = true
							player.Bingo = 1
							if p.Name == myName || p.Name == myToken {
								myX = j
								myY = i
								myGold = p.Gold
								records[roundId].Gold = p.Gold
								records[roundId].RoundId = roundId
								records[roundId].X = myX
								records[roundId].Y = myY
								records[roundId].Crowded = playerLen
							} else if p.Name == "FixedRobot" {
								fixedX = j
								fixedY = i
								fixedGold = p.Gold
								player.Precision = 1 // 固定机器人一定是准的...
							} else if p.Name == "RandRobot" {
								player.Precision = 0 // 随机机器人一定不准
							} else {
								player.Precision = 1
							}
							players = append(players, player)
						} else {
							player = &PlayerInfo{}
							player.X = j
							player.Y = i
							player.Name = p.Name
							player.Gold = p.Gold
							player.Joined = true
							if p.Name == myName || p.Name == myToken {
								myX = j
								myY = i
								myGold = p.Gold
								records[roundId].Gold = p.Gold
								records[roundId].RoundId = roundId
								records[roundId].X = myX
								records[roundId].Y = myY
								records[roundId].Crowded = playerLen
							} else if p.Name == "FixedRobot" {
								fixedX = j
								fixedY = i
								fixedGold = p.Gold
								player.Precision = 1 // 固定机器人一定是准的...
							} else if p.Name == "RandRobot" {
								player.Precision = 0 // 随机机器人一定不准
							} else {
								player.HighestPrecision = 0
								player.LowestPrecision = 1
								player.Precision = 1
							}
							allPlayers[p.Name] = player
							players = append(players, player)
						}
					}
				}
			}
		}
	} else {
		// 非第一局, 理论上不用再处理新玩家了, 如果在map中找不到, 就不管了
		for i := 0; i < len(tileMap); i++ {
			for j := 0; j < len(tileMap[i]); j++ {
				tileMap[i][j].x = j
				tileMap[i][j].y = i
				if len(tileMap[i][j].Players) > 0 {
					playerLen := len(tileMap[i][j].Players)
					for k := 0; k < playerLen; k++ {
						p := tileMap[i][j].Players[k]
						player, ok := allPlayers[p.Name]
						if ok {
							player.X = j
							player.Y = i
							player.Gold = p.Gold
							if p.Name == myName || p.Name == myToken {
								myX = j
								myY = i
								myGold = p.Gold
								records[roundId].Gold = p.Gold
								records[roundId].RoundId = roundId
								records[roundId].X = myX
								records[roundId].Y = myY
								records[roundId].Crowded = playerLen
							} else if p.Name == "FixedRobot" {
								fixedX = j
								fixedY = i
								fixedGold = p.Gold
							} else if p.Name != "RandRobot" {
								for l := 0; l < len(player.Cells); l++ {
									if player.X == player.Cells[l].X && player.Y == player.Cells[l].Y {
										player.Bingo++
										break
									}
								}
								player.Precision = float32(player.Bingo) / float32(roundId)
							}
							players = append(players, player)
						}
					}
				}
			}
		}
	}

	playersLen := len(players)
	// 给所有玩家排个序
	for i := 0; i < playersLen; i++ {
		for j := i + 1; j < playersLen; j++ {
			if players[i].Gold < players[j].Gold {
				temp := players[i]
				players[i] = players[j]
				players[j] = temp
			}
		}
	}

	records[roundId].FirstGold = players[0].Gold
	records[roundId].FirstName = players[0].Name

	var myRank float32 = 1
	for i := 0; i < playersLen; i++ {
		if players[i].Name == myName || players[i].Name == myToken {
			myRank = float32(i+1) / float32(playersLen)
			records[roundId].Rank = myRank
			break
		}
	}

	// 先判断固定机器人是否金币比我多
	if fixedGold > myGold && (fixedGold <= 0 || fixedGold%5 != 0 || myRank >= 0.2) {
		fixedCost := int(math.Floor((math.Abs(float64(fixedX-myX)) + math.Abs(float64(fixedY-myY))) * 1.5))
		if fixedCost <= myGold {
			records[roundId].TargetX = fixedX
			records[roundId].TargetY = fixedY
			records[roundId].Expected = 1000 + float32(fixedGold)*0.25 // 特别标记
			sendMove(fixedX, fixedY, roundId)
			fmt.Println(records[roundId])
			return
		}

		tileMap[fixedY][fixedX].cost = fixedCost
		tileMap[fixedY][fixedX].left = myGold - fixedCost
		tileMap[fixedY][fixedX].leftFlag = 1
		tileMap[fixedY][fixedX].risk = 10
	}

	// 预测所有玩家行为
	for i := 0; i < playersLen; i++ {
		player := players[i]
		player.Rank = float32(i+1) / float32(playersLen)
		if player.Name == myName || player.Name == myToken || player.Name == "FixedRobot" || player.Name == "RandRobot" {
			continue
		}
		player.Cells = make([]*Cell, 0, 4)
		var highestExp int = -99
		var left int = 0
		for j := 0; j < len(tileMap); j++ {
			for k := 0; k < len(tileMap[j]); k++ {
				var cost int
				if k == player.X && j == player.Y {
					cost = 1
				} else {
					cost = int(math.Floor((math.Abs(float64(k-player.X)) + math.Abs(float64(j-player.Y))) * 1.5))
				}
				expected := tileMap[j][k].Gold - cost
				if expected == highestExp {
					player.Cells = append(player.Cells, &Cell{X: k, Y: j})
				}
				if expected > highestExp {
					player.Cells = player.Cells[0:0] // 清空切片
					player.Cells = append(player.Cells, &Cell{X: k, Y: j})
					highestExp = expected
					left = player.Gold - cost
				}
			}
		}
		player.EstimateLeft = left
		precision := (player.Precision + player.AverPrecision) / 2
		for j := 0; j < len(player.Cells); j++ {
			cell := player.Cells[j]
			tileMap[cell.Y][cell.X].reachable += 1
			var myLeft int = tileMap[cell.Y][cell.X].left
			if tileMap[cell.Y][cell.X].leftFlag == 0 {
				var cost int
				if cell.X == myX && cell.Y == myY {
					cost = 1
				} else {
					cost = int(math.Floor((math.Abs(float64(cell.X-myX)) + math.Abs(float64(cell.Y-myY))) * 1.5))
				}
				myLeft = myGold - cost
				tileMap[cell.Y][cell.X].left = myLeft
				tileMap[cell.Y][cell.X].cost = cost
				tileMap[cell.Y][cell.X].leftFlag = 1
			}

			if left > myLeft {
				tileMap[cell.Y][cell.X].opportunity += precision
			} else {
				tileMap[cell.Y][cell.X].risk += precision
			}
		}
	}

	// 根据棋盘数据生成 cells info
	cells := make([]*Tile, 0, 48)
	for i := 0; i < len(tileMap); i++ {
		for j := 0; j < len(tileMap[i]); j++ {
			tile := tileMap[i][j]
			if tile.leftFlag == 1 && tile.left < 0 {
				continue
			} else if tile.leftFlag == 0 {
				var cost int
				if j == myX && i == myY {
					cost = 1
				} else {
					cost = int(math.Floor((math.Abs(float64(tile.x-myX)) + math.Abs(float64(tile.y-myY))) * 1.5))
				}
				left := myGold - cost
				if left < 0 {
					continue
				}
				tile.leftFlag = 1
				tile.left = left
				tile.cost = cost
			}
			cells = append(cells, tileMap[i][j])

			// 计算预期收益
			if myRank < 0.1 || myGold >= lastGold {
				// 第一名的时候, 每次都找负分数的地方躲
				if tile.Gold == 0 || tile.Gold == 1 || tile.Gold == -1 {
					tile.expected = 15 - tile.risk - float32(tile.cost)
				} else {
					tile.expected = float32(tile.Gold) - tile.risk*2 - float32(tile.cost)
				}
			} else if myRank > 0.5 {
				// 倒数, 考虑一个激进的策略
				if tile.Gold > 0 && tile.Gold%5 == 0 {
					tile.expected = float32(tile.Gold) + tile.opportunity*2 - float32(tile.cost)
				} else {
					tile.expected = float32(tile.Gold) + tile.opportunity - tile.risk - float32(tile.cost)
				}
			} else {
				tile.expected = float32(tile.Gold) + tile.opportunity - tile.risk - float32(tile.cost)
			}
		}
	}

	// cells 根据 expected 来排序
	for i := 0; i < len(cells); i++ {
		for j := i + 1; j < len(cells); j++ {
			if cells[i].expected < cells[j].expected {
				temp := cells[i]
				cells[i] = cells[j]
				cells[j] = temp
			}
		}
	}

	// 找出预期值最高的n个, 再随机选一个
	if len(cells) == 0 {
		records[roundId].TargetX = myX
		records[roundId].TargetY = myY
		records[roundId].Expected = 2000 + float32(tileMap[myY][myX].Gold)
		sendMove(myX, myY, frameData.RoundID)
		fmt.Println(records[roundId])
		return
	}

	highest := make([]*Tile, 0, 8)
	bestEx := cells[0].expected
	highest = append(highest, cells[0])
	for i := 0; i < len(cells); i++ {
		if cells[i].expected == bestEx {
			highest = append(highest, cells[i])
		} else {
			break
		}
	}

	// 做个随机, 让自己的行为难以预测
	lenOfHighest := int32(len(highest))
	rnd := rand.Int31n(lenOfHighest)
	if rnd >= lenOfHighest {
		rnd = lenOfHighest - 1
	}

	records[roundId].TargetX = cells[rnd].x
	records[roundId].TargetY = cells[rnd].y
	records[roundId].Expected = cells[rnd].expected
	sendMove(cells[rnd].x, cells[rnd].y, frameData.RoundID)
	fmt.Println(records[roundId])
}

var myX, myY, myGold int
var fixedX, fixedY, fixedGold int
var lastGold int = 100

var myName string = "追光骑士团"
var myToken string = "DyKiQSgpDhrtMQSsVgvs7NWtS7A79XLI"

var ws *websocket.Conn
var records [96]Record

var allPlayers map[string]*PlayerInfo = make(map[string]*PlayerInfo)

func main() {
	//online
	uri := "ws://pgame.51wnl-cq.com:8881/ws"
	//local debug
	//uri := "ws://localhost:8881/ws"
	rand.Seed(time.Now().UnixNano())
LOGIN:
	err := login(uri, myToken)
	if err != nil {
		goto LOGIN
	}
	defer func() {
		if ws != nil {
			ws.Close()
		}
	}()

PREPARE:
	for {
		time.Sleep(time.Millisecond * 10) // 睡 10 毫秒
		data := Msg{}
		err = recvMessage(&data)
		if err == nil {
			if data.Msgtype == 1 {
				fmt.Println("准备")
				s := Msg{}
				s.Msgtype = 2
				s.Token = myToken
				sendMessage(s)
				goto GAMELOOP
			}
		} else {
			goto LOGIN
		}
	}
GAMELOOP:
	for {
		time.Sleep(time.Millisecond * 10) // 睡 10 毫秒
		data := Game{}
		err = recvMessage(&data)
		if err == nil {
			if data.Msgtype == 5 {
				saveGame(data.GameID)
				lastGold = myGold
				fmt.Println("游戏结束, GameId = ", data.GameID)
				for i := 0; i < len(data.Sorted); i++ {
					fmt.Println("Name = ", data.Sorted[i].Name, ", Gold = ", data.Sorted[i].Gold)
				}
				//fmt.Println(records)
				goto PREPARE
			} else if data.Msgtype == 3 {
				updateFrame(data)
			}
		} else {
			fmt.Println(err)
			goto LOGIN
		}
	}
}
