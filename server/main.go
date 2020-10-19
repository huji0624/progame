package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func LogDebug(args ...interface{}) {
	log.Println(args...)
}

func WhateverOrigin(r *http.Request) bool {
	return true
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

type PlayerInfo struct {
	Key  string
	X    int
	Y    int
	Gold int
}

type Player struct {
	c  *websocket.Conn
	rc chan *Msg

	token string
	Info  *PlayerInfo
}

var upgrader = websocket.Upgrader{CheckOrigin: WhateverOrigin} // use default options
var connections = make(map[string]*Player)

var prepares = make(map[string]*Player)

func LogStruct(v interface{}) {
	bt, err := json.MarshalIndent(v, "", "\t")
	if err == nil {
		LogDebug(string(bt))
	}
}

func WriteToClient(c *websocket.Conn, v interface{}) {
	// log.Println("will write data:")
	// LogStruct(v)
	c.WriteJSON(v)
}

func WriteToClientData(c *websocket.Conn, data []byte) {
	err := c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write err:", err)
	} else {
		// LogDebug("write data:")
		// LogStruct(string(data))
	}
}

func WaitForMsg(c *websocket.Conn, v interface{}) error {
	err := c.ReadJSON(v)
	return err
}

var mutex sync.Mutex

func echo(w http.ResponseWriter, r *http.Request) {
	log.Println("new connection..")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// outtime := time.Now().Add(1 * time.Second)
	// c.SetReadDeadline(outtime)

	jmsg := Msg{}
	err = WaitForMsg(c, &jmsg)
	if err != nil {
		log.Println("read err:", err)
	} else {
		if jmsg.Msgtype == 0 {
			allowed := false
			name := ""
			ok := false
			if tokenmap == nil {
				name = jmsg.Token
				ok = true
			} else {
				name, ok = tokenmap[jmsg.Token]
			}

			mutex.Lock()
			defer mutex.Unlock()

			if ok {
				_, online := prepares[name]
				if !online {
					allowed = true
				} else {
					allowed = false
				}
			} else {
				allowed = false
			}

			if allowed {
				rc := make(chan *Msg, 200)
				p := &Player{c: c, rc: rc}
				p.token = name
				prepares[name] = p
				WriteToClient(c, &jmsg)
				go ReadLoop(p, jmsg.Token)
			} else {
				jmsg.Msgtype = -1
				WriteToClient(c, &jmsg)
				c.Close()
			}

		} else {
			jmsg.Msgtype = -1
			WriteToClient(c, &jmsg)
			c.Close()
		}
	}
}

func KickPlayer(p *Player) {
	log.Println("will kick player:", p.token)
	p.c.Close()
	delete(connections, p.token)
	delete(prepares, p.token)
	log.Println("did kick player:", p.token)
}

const (
	MapWidth  = 8
	MapHeight = 6
)

type Game struct {
	GameID  uint64
	Msgtype int
	status  int //-1无效0准备1开始
	RoundID int
	Wid     int
	Hei     int

	Tilemap [MapHeight][MapWidth]*Tile

	roundRecords []string
}

func initGame(g *Game, playings map[string]*Player) {
	log.Println("init game.")
	g.RoundID = 0
	g.GameID = g.GameID + 1
	g.roundRecords = make([]string, 0, 0)
	g.Wid = MapWidth
	g.Hei = MapHeight
	for i := 0; i < MapWidth; i++ {
		for j := 0; j < MapHeight; j++ {
			t := &Tile{Gold: 0}
			t.players = make(map[string]*Player, 0)
			g.Tilemap[j][i] = t
		}
	}

	if records != nil && records.Moves == nil {
		records.Moves = make(map[string]uint64)
	}
	if records != nil && records.Gives == nil {
		records.Gives = make(map[string]uint64)
	}

	for token, v := range playings {
		x := rand.Intn(MapWidth)
		y := rand.Intn(MapHeight)
		v.Info = &PlayerInfo{X: x, Y: y}
		v.Info.Gold = 0
		v.Info.Key = token

		MovePlayerForce(g, v, x, y)
	}
}

func PublicToClientData(data []byte) {
	for _, v := range connections {
		if v.rc != nil {
			WriteToClientData(v.c, data)
		}
	}
}

func pubGameMap(g *Game) []byte {
	for i := 0; i < MapWidth; i++ {
		for j := 0; j < MapHeight; j++ {
			t := g.Tilemap[j][i]
			if len(t.players) > 0 {
				tmp := make([]*GameScore, 0, 3)
				for _, v := range t.players {
					tmp = append(tmp, &GameScore{Name: v.Info.Key, Gold: v.Info.Gold})
				}

				t.P = tmp
			} else {
				t.P = nil
			}
		}
	}

	// LogDebug("will pub data:")
	// LogStruct(g)

	jmsg, err := json.Marshal(g)
	if err != nil {
		log.Println("pubGame err:", err)
	} else {
		PublicToClientData(jmsg)
	}

	return jmsg
}

func pubGameResults() {
	tmsg := &Msg{Msgtype: 5, Sorted: records.Sorted}

	jmsg, err := json.Marshal(tmsg)
	if err != nil {
		log.Println("pubGame err:", err)
	} else {
		PublicToClientData(jmsg)
	}
}

func AbsI(n int) int {
	if n < 0 {
		return -n
	} else {
		return n
	}
}

func MaxI(n int, m int) int {
	if n < m {
		return m
	} else {
		return n
	}
}

func MovePlayerForce(g *Game, player *Player, x int, y int) {
	info := player.Info

	t := g.Tilemap[info.Y][info.X]

	info.X = x
	info.Y = y

	tt := g.Tilemap[info.Y][info.X]

	delete(t.players, player.Info.Key)
	tt.players[info.Key] = player
}

func MovePlayer(g *Game, player *Player, x int, y int) {
	info := player.Info

	//player move cost gold.And get one gold if not move.
	step := AbsI(x-info.X) + AbsI(y-info.Y)

	//check if move out of the map.dont move at all.
	if x < 0 || y < 0 || x >= g.Wid || y >= g.Hei {
		step = 0
	}

	cost := int(float64(step) * 1.5)

	if step == 0 {
		// log.Println("Player not moving!", player.token, x, ",", y)
		player.Info.Gold--
		return
	} else if cost <= info.Gold {
		player.Info.Gold = player.Info.Gold - cost
	} else {
		player.Info.Gold--
		// log.Println("Gold not enough!", player.token)
		return
	}

	if records != nil {
		records.Moves[player.Info.Key] = records.Moves[player.Info.Key] + uint64(step)
	}

	MovePlayerForce(g, player, x, y)
}

func CheckGameOver(g *Game) bool {
	if g.RoundID == 96 {
		return true
	}

	return false
}

type Tile struct {
	Gold    int
	P       []*GameScore `json:"Players,omitempty"`
	players map[string]*Player
}

func GetRandomPlayer(players map[string]*Player) *Player {
	tmp := make([]*Player, 0, 3)
	for _, v := range players {
		tmp = append(tmp, v)
	}

	index := rand.Intn(len(tmp))
	p := tmp[index]
	return p
}

func ApplyGameLogic(g *Game, playings map[string]*Player) {

	for i := 0; i < MapWidth; i++ {
		for j := 0; j < MapHeight; j++ {
			t := g.Tilemap[j][i]

			//randomly give gold to one player
			if t.Gold != 0 && len(t.players) > 0 {
				if t.Gold == -4 {
					//40% interest or lose 4
					p := GetRandomPlayer(t.players)
					r := rand.Intn(2)
					if r == 0 {
						p.Info.Gold = p.Info.Gold + p.Info.Gold/10*4
					} else {
						p.Info.Gold = p.Info.Gold + t.Gold
					}
				} else if t.Gold == 7 || t.Gold == 11 {
					//1 player get.multi lose.
					pcount := len(t.players)
					p := GetRandomPlayer(t.players)
					if pcount == 1 {
						p.Info.Gold += t.Gold
					} else {
						p.Info.Gold -= t.Gold
					}
				} else if t.Gold > 0 && t.Gold%5 == 0 {
					pcount := len(t.players)
					if pcount == 1 {
						p := GetRandomPlayer(t.players)
						p.Info.Gold += t.Gold
					} else {
						total := t.Gold
						for _, v := range t.players {
							total += v.Info.Gold
						}
						each := total / len(t.players)
						for _, v := range t.players {
							v.Info.Gold = each
						}
					}
				} else if t.Gold == 8 {
					//share Gold
					pcount := len(t.players)
					each := t.Gold / pcount
					for _, v := range t.players {
						v.Info.Gold += each
					}
				} else {
					p := GetRandomPlayer(t.players)
					p.Info.Gold += t.Gold
				}

				t.Gold = 0
			}

			//who has more gold will give 1/4 of his gold to others
			if len(t.players) > 1 {
				GiveGoldProcess(t)
			}
		}
	}
}

func HalfPlayersGold(behalfs map[string]*Player) {
	for _, v := range behalfs {
		v.Info.Gold = v.Info.Gold / 2
	}
}

func GiveGoldProcess(t *Tile) {
	mostGold := -1
	mostGoldCount := 0
	for _, v := range t.players {
		// log.Println("player : ", v.Info.Key, " = ", v.Info.Gold)
		if v.Info.Gold > mostGold {
			mostGold = v.Info.Gold
			mostGoldCount = 1
		} else if v.Info.Gold == mostGold {
			mostGoldCount++
		}
	}
	poor := len(t.players) - mostGoldCount
	if poor > 0 {
		give := mostGold / 4
		each := give / poor
		if each > 0 {
			for _, v := range t.players {
				if v.Info.Gold == mostGold {
					v.Info.Gold = v.Info.Gold - poor*each
					if records != nil {
						records.Gives[v.Info.Key] = records.Gives[v.Info.Key] + uint64(poor*each)
					}
				} else {
					v.Info.Gold = v.Info.Gold + mostGoldCount*each
				}
			}
		}
	}
}

func RandomGenGold(g *Game, playings map[string]*Player) {

	n := MapHeight * MapWidth / 6
	for i := 0; i < n; i++ {
		r1 := rand.Intn(8)
		r2 := rand.Intn(6)
		r := r1 - r2

		x := rand.Intn(MapWidth)
		y := rand.Intn(MapHeight)

		generateAt(g, x, y, r)
	}
}

func generateAt(g *Game, x int, y int, n int) {
	t := g.Tilemap[y][x]
	t.Gold += n
}

func ReadLoop(p *Player, token string) {
	for {
		// if p.playing {
		// 	outtime := time.Now().Add(1 * time.Second)
		// 	p.c.SetReadDeadline(outtime)
		// }

		jmsg := Msg{}
		err := p.c.ReadJSON(&jmsg)
		// log.Println("read loop:")
		// LogStruct(jmsg)
		if err == nil {
			if p.rc != nil {
				p.rc <- &jmsg
			} else {
				log.Println("channel nil.wait next...")
			}
		} else {
			log.Println("ReadLoop read token:", token, " err:", err, " player will not move.")
			close(p.rc)
			p.rc = nil
			break
		}
	}
}

func RunPlayerMoves(game *Game, playings map[string]*Player) {
	//play loop
	gmsg := pubGameMap(game)
	game.roundRecords = append(game.roundRecords, string(gmsg))

	time.Sleep(time.Second)

	//wait for move
	for _, v := range playings {
		var msg *Msg
		for {
			conti := false
			select {
			case msg = <-v.rc:
				if msg.RoundID == game.RoundID {
					conti = false
				} else {
					conti = true
				}
			default:
				msg = &Msg{}
				msg.Msgtype = -1
				conti = false
			}
			if !conti {
				break
			}
		}

		if msg.Msgtype == 4 {
			MovePlayer(game, v, msg.X, msg.Y)
		} else if msg.Msgtype == -1 {
			MovePlayer(game, v, v.Info.X, v.Info.Y)
		} else {
			KickPlayer(v)
		}
	}
}

func PlayOneRound(game *Game, playings map[string]*Player, playermove func(game *Game, playings map[string]*Player), randomGold bool) {
	log.Println("Will Start Game(", game.GameID, ") Round(", game.RoundID, ")..")

	//1 gold each round.
	for _, v := range playings {
		v.Info.Gold++
	}

	//give out interest
	for _, v := range playings {
		v.Info.Gold = v.Info.Gold + v.Info.Gold/10*2
	}

	if randomGold {
		//random Gold
		RandomGenGold(game, playings)
	}

	playermove(game, playings)

	game.RoundID++
	ApplyGameLogic(game, playings)
}

func PlayGameRounds(game *Game, playings map[string]*Player) {
	initGame(game, connections)
	for {
		if len(prepares) == 0 {
			break
		}

		if CheckGameOver(game) {
			game.status = -1
			log.Println("Game Over.")
			SaveGameResult(game)
			pubGameResults()
			break
		} else {
			PlayOneRound(game, playings, RunPlayerMoves, true)
		}
	}
}

func GameLoop() {
	game := &Game{status: -1, Msgtype: 3}
	g_game = game
	con, err := ioutil.ReadFile(recordsPath + "/gameid")
	if err == nil {
		game.GameID, _ = strconv.ParseUint(string(con), 0, 64)
	} else {
		log.Println("read game id err:", err)
	}
	log.Println("Game ID:", game.GameID)

	for {
		log.Println("will prepare for next game.Game id:", game.GameID)
		// if game.GameID == 1 {
		// 	time.Sleep(time.Hour * 12)
		// }
		for {
			if len(prepares) <= 1 {
				time.Sleep(time.Second)
				log.Println("no player.go on...")
			} else {
				break
			}
		}

		game.status = 0
		log.Println("prepare for next game.")
		//recreate channel
		for _, v := range prepares {
			trc := v.rc
			v.rc = nil
			if trc != nil {
				close(trc)
			}
			v.rc = make(chan *Msg, 200)
		}

		//prepare
		jmsg := Msg{}
		jmsg.Msgtype = 1
		for _, v := range prepares {
			jmsg.Token = v.token
			WriteToClient(v.c, jmsg)
		}

		time.Sleep(time.Second * 3)

		log.Println("wait for response.")
		connections = make(map[string]*Player)
		//prepare response
		jmsg = Msg{}
		for token, v := range prepares {
			log.Println("check ready:", token)

			select {
			case msg, ok := <-v.rc:
				if ok && msg.Msgtype == 2 {
					log.Println("player is ready:", token)

					connections[v.token] = v
				} else {
					KickPlayer(v)
				}
			default:
				log.Println("no prepare response from:", token)
				KickPlayer(v)
			}
		}

		if len(connections) == 0 {
			continue
		}

		log.Println("will play game rounds.players:", len(connections))
		game.status = 1
		PlayGameRounds(game, connections)
	}
}

func SaveGameResult(g *Game) {
	if records.Scores == nil {
		records.Scores = make(map[int]map[string]int)
	}

	tmpScore := make(map[string]int)
	for _, p := range connections {
		tmpScore[p.Info.Key] = p.Info.Gold
	}
	records.Scores[records.Index] = tmpScore
	records.Index++
	if records.Index == 96 {
		records.Index = 0
	}

	totalScore := make(map[string]int)
	totalCount := make(map[string]int)
	ps := make([]*GameScore, 0, len(connections))
	for _, v := range records.Scores {
		for key, score := range v {
			totalScore[key] = totalScore[key] + score
			totalCount[key] = totalCount[key] + 1
		}
	}
	for k, v := range totalScore {
		av := v / totalCount[k]
		ps = append(ps, &GameScore{Name: k, Gold: av})
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Gold > ps[j].Gold })

	records.Sorted = ps

	LogStruct(records)
	grecord, _ := json.Marshal(g.roundRecords)
	fp := fmt.Sprintf("%v/game_%v.json", recordsPath, g.GameID)
	if err := ioutil.WriteFile(fp, grecord, 0644); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("write file success...", fp)
		ioutil.WriteFile(recordsPath+"/gameid", []byte(fmt.Sprintf("%v", g.GameID)), 0644)
	}
}

type GameRank struct {
	// Gameresults []
	Sorted []*GameScore

	Index  int
	Scores map[int]map[string]int

	Moves map[string]uint64
	Gives map[string]uint64
}

type GameScore struct {
	Name string
	Gold int
}

type Rank struct {
	Total  []*GameScore
	First  []*GameScore
	Second []*GameScore
	Third  []*GameScore
	Gid    uint64
}

func setupresponse(w http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		log.Println("origin is :", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	} else {
		log.Println("origin is nil.")
		log.Println(origin)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	setupresponse(w, r)

	r.ParseForm()
	gid := r.Form.Get("gid")

	bt, err := ioutil.ReadFile(recordsPath + "/game_" + gid + ".json")
	if err == nil {
		_, _ = w.Write(bt)
	} else {
		_, _ = w.Write([]byte("no such game."))
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	setupresponse(w, r)

	var err string

	r.ParseForm()
	token := r.Form.Get("token")
	name := r.Form.Get("name")
	switch name {
	case "first":
	case "second":
	case "third":
	default:
		err = "params err"
	}
	if token != "iaijmnxuahiqooqjs8918221h" {
		err = "token fail."
	}
	if err != "" {
		_, _ = w.Write([]byte(err))
		return
	}

	bt, _ := json.Marshal(records)
	reterr := ioutil.WriteFile(recordsPath+"/"+name+".json", bt, 0644)

	if reterr != nil {
		_, _ = w.Write([]byte(reterr.Error()))
	} else {
		_, _ = w.Write([]byte(name + " DONE"))
	}
}

func rankHandler(w http.ResponseWriter, r *http.Request) {
	setupresponse(w, r)

	ret := &Rank{}

	pathes := [3]string{"first.json", "second.json", "third.json"}

	for _, v := range pathes {
		tpath := recordsPath + "/" + v
		if r, _ := exists(tpath); r {
			result := &GameRank{}
			b, err := ioutil.ReadFile(tpath)
			if err != nil {
				log.Println("err read file:", err)
				continue
			}
			err = json.Unmarshal(b, result)
			if err != nil {
				log.Println("err read rank.")
				return
			} else if v == "first.json" {
				ret.First = result.Sorted
			} else if v == "second.json" {
				ret.Second = result.Sorted
			} else if v == "third.json" {
				ret.Third = result.Sorted
			}
		}
	}

	ret.Total = records.Sorted

	ret.Gid = g_game.GameID

	bt, err := json.Marshal(ret)
	if err != nil {
		log.Println("err marsh json rank.")
		return
	} else {
		_, _ = w.Write(bt)
	}
}

var records *GameRank
var recordsPath string
var g_game *Game
var tokenmap map[string]string

func main() {
	rand.Seed(time.Now().UnixNano())

	if !RunUnitTest() {
		log.Println("Unit test fail...exit")
		os.Exit(1)
		return
	}

	pid := os.Getpid()
	s := fmt.Sprintln(pid)
	ioutil.WriteFile("./pid.log", []byte(s), 0655)

	recordsPath = "records"

	os.Mkdir(recordsPath, os.ModePerm)

	records = &GameRank{}

	tpath := recordsPath + "/tokenmap.json"
	b, ex := ioutil.ReadFile(tpath)
	if ex == nil {
		json.Unmarshal(b, &tokenmap)
		LogStruct(tokenmap)
	} else {
		tokenmap = nil
	}

	go GameLoop()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/rank", rankHandler)
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/ws", echo)
	log.Println("ws server ready...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8881", nil))
}

//===========Unit test==============

func RunUnitTest() bool {
	testGame := &Game{}
	testplayers := make(map[string]*Player)

	a := &Player{}
	b := &Player{}
	c := &Player{}
	d := &Player{}

	testplayers["a"] = a
	testplayers["b"] = b
	testplayers["c"] = c
	testplayers["d"] = d

	initGame(testGame, testplayers)

	//test all moves
	MovePlayerForce(testGame, a, 0, 0)
	MovePlayerForce(testGame, b, 1, 1)
	MovePlayerForce(testGame, c, 7, 4)
	MovePlayerForce(testGame, d, 2, 5)

	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {

		MovePlayer(testGame, a, 0, 0)
		MovePlayer(testGame, b, 1, 2)
		MovePlayer(testGame, c, 8, 4)
		MovePlayer(testGame, d, 7, 5)

	}, false)

	if a.Info.X == 0 && a.Info.Y == 0 && a.Info.Gold == 0 {

	} else {
		log.Println("move a to same position test fail.")
		return false
	}

	if b.Info.X == 1 && b.Info.Y == 2 && b.Info.Gold == 0 {

	} else {
		log.Println("move b to 1,2 test fail.")
		return false
	}

	if c.Info.X == 7 && c.Info.Y == 4 && b.Info.Gold == 0 {

	} else {
		log.Println("move c to out range test fail.")
		return false
	}

	if d.Info.X == 2 && d.Info.Y == 5 && b.Info.Gold == 0 {

	} else {
		log.Println("move d not enough gold test fail.")
		return false
	}

	a.Info.Gold = 9
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {

		MovePlayer(testGame, a, 3, 3)

	}, false)

	if a.Info.X == 3 && a.Info.Y == 3 && a.Info.Gold == 3 {

	} else {
		log.Println("move a to 3,3 position test fail.")
		LogStruct(testplayers)
		return false
	}
	//a 2 (3,3)

	a.Info.Gold = 2
	//test player get gold
	generateAt(testGame, 3, 4, 3)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {

		MovePlayer(testGame, a, 3, 4)

	}, false)
	if a.Info.X == 3 && a.Info.Y == 4 && a.Info.Gold == 5 {

	} else {
		log.Println("get gold test fail.")
		LogStruct(testplayers)
		return false
	}

	generateAt(testGame, 3, 4, -1)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {

		MovePlayer(testGame, a, 3, 4)

	}, false)
	if a.Info.X == 3 && a.Info.Y == 4 && a.Info.Gold == 4 {

	} else {
		log.Println("get -1 gold test fail.")
		LogStruct(testplayers)
		return false
	}

	a.Info.Gold = 2
	d.Info.Gold = 2
	generateAt(testGame, 3, 5, 1)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {

		MovePlayer(testGame, a, 3, 5)
		MovePlayer(testGame, d, 3, 5)

	}, false)

	if a.Info.X == 3 && a.Info.Y == 5 && d.Info.X == 3 && d.Info.Y == 5 && ((a.Info.Gold == 2 && d.Info.Gold == 3) || (a.Info.Gold == 3 && d.Info.Gold == 2)) {

	} else {
		log.Println("2 players get gold test fail.")
		LogStruct(testplayers)
		return false
	}

	//test player give gold
	a.Info.Gold = 8
	d.Info.Gold = 2

	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 3, 5)
		MovePlayer(testGame, d, 3, 5)
	}, false)

	if a.Info.X == 3 && a.Info.Y == 5 && d.Info.X == 3 && d.Info.Y == 5 && a.Info.Gold == 6 && d.Info.Gold == 4 {

	} else {
		log.Println("player gives gold test fail.")
		LogStruct(testplayers)
		return false
	}

	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 3, 5)
		MovePlayer(testGame, d, 3, 5)
	}, false)

	if a.Info.X == 3 && a.Info.Y == 5 && d.Info.X == 3 && d.Info.Y == 5 && a.Info.Gold == 5 && d.Info.Gold == 5 {

	} else {
		log.Println("player gives gold test fail.")
		LogStruct(testplayers)
		return false
	}

	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 3, 5)
		MovePlayer(testGame, d, 3, 5)
	}, false)

	if a.Info.X == 3 && a.Info.Y == 5 && d.Info.X == 3 && d.Info.Y == 5 && a.Info.Gold == 5 && d.Info.Gold == 5 {

	} else {
		log.Println("player gives gold test fail.")
		LogStruct(testplayers)
		return false
	}

	//test gold get interest
	c.Info.Gold = 20
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, c, c.Info.X, c.Info.Y)
	}, false)

	if c.Info.X == 7 && c.Info.Y == 4 && c.Info.Gold == 24 {

	} else {
		log.Println("player get interest test fail.")
		return false
	}

	//test magic gold -4
	c.Info.Gold = 10
	generateAt(testGame, c.Info.X, c.Info.Y, -4)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, c, c.Info.X, c.Info.Y)
	}, false)

	if c.Info.X == 7 && c.Info.Y == 4 && (c.Info.Gold == 16 || c.Info.Gold == 8) {

	} else {
		log.Println("magic gold -4 test fail.")
		LogStruct(c)
		return false
	}

	//test magic gold 5
	c.Info.Gold = 0
	generateAt(testGame, c.Info.X, c.Info.Y, 5)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, c, c.Info.X, c.Info.Y)
	}, false)

	if c.Info.X == 7 && c.Info.Y == 4 && c.Info.Gold == 5 {

	} else {
		log.Println("magic gold 5 test fail.")
		return false
	}

	//test magic gold 10
	a.Info.Gold = 1
	d.Info.Gold = 3
	generateAt(testGame, 3, 5, 10)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 3, 5)
		MovePlayer(testGame, d, 3, 5)
	}, false)

	if a.Info.X == 3 && a.Info.Y == 5 && d.Info.X == 3 && d.Info.Y == 5 && a.Info.Gold == 7 && d.Info.Gold == 7 {

	} else {
		log.Println("a b magic gold 10 test fail.")
		LogStruct(testplayers)
		return false
	}

	//test magic gold  8
	a.Info.Gold = 2
	d.Info.Gold = 4
	generateAt(testGame, 2, 5, 8)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 2, 5)
		MovePlayer(testGame, d, 2, 5)
	}, false)

	if a.Info.X == 2 && a.Info.Y == 5 && d.Info.X == 2 && d.Info.Y == 5 && a.Info.Gold == 6 && d.Info.Gold == 8 {

	} else {
		log.Println("a b magic gold 8 test fail.")
		LogStruct(testplayers)
		return false
	}

	//test magic gold  7
	a.Info.Gold = 2
	generateAt(testGame, 4, 5, 7)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 4, 5)
	}, false)

	if a.Info.X == 4 && a.Info.Y == 5 && a.Info.Gold == 9 {

	} else {
		log.Println("a  magic gold 7 test fail.")
		LogStruct(testplayers)
		return false
	}

	//test magic gold  11
	a.Info.Gold = 12
	d.Info.Gold = 14
	generateAt(testGame, 2, 4, 11)
	PlayOneRound(testGame, testplayers, func(game *Game, playings map[string]*Player) {
		MovePlayer(testGame, a, 2, 4)
		MovePlayer(testGame, d, 2, 4)
	}, false)

	if a.Info.X == 2 && a.Info.Y == 4 && d.Info.X == 2 && d.Info.Y == 4 && (a.Info.Gold == 1 || d.Info.Gold == 3) {

	} else {
		log.Println("a b magic gold 11 test fail.")
		LogStruct(testplayers)
		return false
	}

	LogStruct(testGame)
	LogStruct(testplayers)

	return true
}
