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
	msgtype 0登陆 -1登陆失败 1请准备 2准备好了 3游戏信息 4玩家行动
*/
type Msg struct {
	Msgtype int
	Token   string
	RoundID int
	X       int
	Y       int
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

	playing bool
	token   string
	Info    *PlayerInfo
}

var upgrader = websocket.Upgrader{CheckOrigin: WhateverOrigin} // use default options
var connections = make(map[string]*Player)

var prepares = make(map[string]*Player)

func LogStruct(v interface{}) {
	bt, err := json.Marshal(v)
	if err == nil {
		LogDebug(string(bt))
	}
}

func WriteToClient(c *websocket.Conn, v interface{}) {
	log.Println("will write data:")
	LogStruct(v)
	c.WriteJSON(v)
}

func WriteToClientData(c *websocket.Conn, data []byte) {
	err := c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write err:", err)
	} else {
		LogDebug("write data:")
		LogStruct(string(data))
	}
}

func WaitForMsg(c *websocket.Conn, v interface{}) error {
	err := c.ReadJSON(v)
	return err
}

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
			rc := make(chan *Msg)
			p := &Player{c: c, rc: rc}
			p.token = jmsg.Token
			prepares[jmsg.Token] = p
			WriteToClient(c, &jmsg)
			go ReadLoop(p, jmsg.Token)
		} else {
			jmsg.Msgtype = -1
			WriteToClient(c, &jmsg)
		}
	}
}

func ReadLoop(p *Player, token string) {
	for {
		if p.playing {
			outtime := time.Now().Add(1 * time.Second)
			p.c.SetReadDeadline(outtime)
		}

		jmsg := Msg{}
		err := p.c.ReadJSON(&jmsg)
		if err == nil {
			p.rc <- &jmsg
		} else {
			log.Println("read token:", token, " err:", err)
			KickPlayer(p)
			break
		}
	}
}

func KickPlayer(p *Player) {
	p.c.Close()
	delete(connections, p.token)
	delete(prepares, p.token)
	log.Println("kick player:", p.token)
}

type Tile struct {
	Gold    int
	Players map[string]*Player
}

const (
	MapSize = 3
)

type Game struct {
	GameID  uint64
	Msgtype int
	status  int //-1无效0准备1开始
	RoundID int

	Tilemap [MapSize][MapSize]*Tile

	roundRecords []string
}

func initGame(g *Game) {
	log.Println("init game.")
	g.RoundID = 0
	g.GameID = g.GameID + 1
	g.roundRecords = make([]string, 0, 0)
	size := MapSize
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t := &Tile{Gold: 0}
			t.Players = make(map[string]*Player, 0)
			g.Tilemap[i][j] = t
		}
	}

	for token, v := range connections {
		x := rand.Intn(MapSize)
		y := rand.Intn(MapSize)
		v.Info = &PlayerInfo{X: x, Y: y}
		v.Info.Gold = 0
		v.Info.Key = token

		MovePlayer(g, v, x, y)
	}
}

func PublicToClientData(data []byte) {
	LogDebug("will pub data:", string(data))
	for _, v := range connections {
		WriteToClientData(v.c, data)
	}
}

func pubGameMap(g *Game) []byte {

	jmsg, err := json.Marshal(g)
	if err != nil {
		log.Println("pubGame err:", err)
	} else {
		PublicToClientData(jmsg)
	}

	return jmsg
}

func MovePlayer(g *Game, player *Player, x int, y int) {
	info := player.Info

	t := g.Tilemap[info.Y][info.X]

	info.X = x
	info.Y = y

	tt := g.Tilemap[info.Y][info.X]

	delete(t.Players, player.Info.Key)
	tt.Players[info.Key] = player
}

func CheckGameOver(g *Game) bool {
	if g.RoundID > 20 {
		return true
	}

	return false
}

func ApplyGameLogic(g *Game) {
	size := MapSize
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t := g.Tilemap[i][j]
			if t.Gold > 0 && len(t.Players) > 0 {
				tmp := make([]*Player, 0, 3)
				for _, v := range t.Players {
					tmp = append(tmp, v)
				}

				index := rand.Intn(len(tmp))
				p := tmp[index]
				p.Info.Gold += t.Gold
				t.Gold = 0
			}
		}
	}
}

func RandomGenGold(g *Game) {
	n := MapSize - 1
	for i := 0; i < n; i++ {
		r := rand.Intn(MapSize)

		x := rand.Intn(MapSize)
		y := rand.Intn(MapSize)

		t := g.Tilemap[y][x]
		t.Gold += r
	}
}

func PlayGameRounds(game *Game) {
	initGame(game)
	for {
		if len(prepares) == 0 {
			break
		}

		//random Gold
		RandomGenGold(game)

		//play loop
		gmsg := pubGameMap(game)
		game.roundRecords = append(game.roundRecords, string(gmsg))

		time.Sleep(time.Second / 4)

		//wait for move
		for token, v := range connections {
			msg, ok := <-v.rc
			if ok && msg.Msgtype == 4 {

				log.Println("player is moving:", token)
				LogStruct(msg)

				MovePlayer(game, v, msg.X, msg.Y)
			} else {
				KickPlayer(v)
			}
		}

		game.RoundID++
		ApplyGameLogic(game)

		if CheckGameOver(game) {
			game.status = -1
			log.Println("Game Over.")
			SaveGameResult(game)
			break
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
		log.Println("will prepare for next game.")
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
		//prepare
		jmsg := Msg{}
		jmsg.Msgtype = 1
		for _, v := range prepares {
			WriteToClient(v.c, jmsg)
		}

		time.Sleep(time.Second)

		log.Println("wait for response.")
		connections = make(map[string]*Player)
		//prepare response
		jmsg = Msg{}
		for token, v := range prepares {
			msg, ok := <-v.rc
			if ok && msg.Msgtype == 2 {
				log.Println("player is ready:", token)
				v.playing = true
				connections[v.token] = v
			} else {
				KickPlayer(v)
			}
		}

		if len(connections) == 0 {
			continue
		}

		game.status = 1
		PlayGameRounds(game)
	}
}

func SaveGameResult(g *Game) {
	for _, p := range connections {
		records.Scores[p.Info.Key] = records.Scores[p.Info.Key] + p.Info.Gold
	}

	ps := make([]*GameScore, 0, len(connections))
	for k, v := range records.Scores {
		ps = append(ps, &GameScore{Name: k, Gold: v})
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Gold > ps[j].Gold })
	records.Sorted = ps
	records.Count++

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

type GameScore struct {
	Name string
	Gold int
}

type GameRank struct {
	// Gameresults []
	Sorted []*GameScore
	Count  int

	Scores map[string]int
}

func home(w http.ResponseWriter, r *http.Request) {

	out := ""

	bt, err := json.Marshal(records)
	if err != nil {
		out = "error"
	} else {
		out = string(bt)
	}

	_, _ = w.Write([]byte(out))
}

type Rank struct {
	Total  []*GameScore
	First  []*GameScore
	Second []*GameScore
	Third  []*GameScore
	Gid    uint64
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
	r.ParseForm()
	gid := r.Form.Get("gid")

	bt, err := ioutil.ReadFile(recordsPath + "/game_" + gid + ".json")
	if err == nil {
		_, _ = w.Write(bt)
	} else {
		_, _ = w.Write([]byte("no such game."))
	}
}

func rankHandler(w http.ResponseWriter, r *http.Request) {

	ret := &Rank{}

	pathes := [4]string{"total.json", "first.json", "second.json", "third.json"}

	for _, v := range pathes {
		tpath := recordsPath + "/" + v
		if r, _ := exists(tpath); r {
			result := &GameRank{}
			b, _ := ioutil.ReadFile(tpath)
			err := json.Unmarshal(b, result)
			if err != nil {
				log.Println("err read rank.")
				return
			} else if v == "total.json" {
				ret.Total = result.Sorted
			} else if v == "first.json" {
				ret.First = result.Sorted
			} else if v == "second.json" {
				ret.Second = result.Sorted
			} else if v == "third.json" {
				ret.Third = result.Sorted
			}
		}
	}

	if ret.First == nil {
		ret.First = records.Sorted
	} else if ret.Second == nil {
		ret.Second = records.Sorted
	} else if ret.Third == nil {
		ret.Third = records.Sorted
	}

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

func main() {
	recordsPath = "records"

	os.Mkdir(recordsPath, os.ModePerm)

	records = &GameRank{}
	records.Scores = make(map[string]int)

	go GameLoop()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/game", gameHandler)
	http.HandleFunc("/rank", rankHandler)
	http.HandleFunc("/", home)
	http.HandleFunc("/ws", echo)
	// http.HandleFunc("/", home)
	log.Println("ws server ready...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
}
