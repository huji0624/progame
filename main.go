package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"
	"sort"

	"github.com/gorilla/websocket"
)

func LogDebug(args... interface{}){
	// log.Println(args...)
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
	Key string
	X    int
	Y    int
	Gold int
}

type Player struct {
	c  *websocket.Conn
	rc chan *Msg

	playing bool
	token string
	Info *PlayerInfo
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
		if p.playing{
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
	Msgtype int
	status  int //-1无效0准备1开始
	RoundID int

	Tilemap [MapSize][MapSize]*Tile
}

func initGame(g *Game) {
	log.Println("init game.")
	g.RoundID = 0
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

		MovePlayer(g,v,x,y)
	}
}

func PublicToClientData(data []byte) {
	LogDebug("will pub data:", string(data))
	for _, v := range connections {
		WriteToClientData(v.c, data)
	}
}

func pubGameMap(g *Game) {

	jmsg, err := json.Marshal(g)
	if err != nil {
		log.Println("pubGame err:", err)
	} else {
		PublicToClientData(jmsg)
	}
}

func MovePlayer(g *Game,player *Player,x int,y int){
	info := player.Info

	t := g.Tilemap[info.Y][info.X]

	info.X = x
	info.Y = y

	tt := g.Tilemap[info.Y][info.X]

	delete(t.Players,player.Info.Key)
	tt.Players[info.Key] = player
}

func CheckGameOver(g *Game) bool {
	if g.RoundID>20{
		return true
	}

	return false
}

func ApplyGameLogic(g *Game){
	size := MapSize
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t :=  g.Tilemap[i][j]
			if t.Gold>0 && len(t.Players)>0{
				tmp := make([]*Player,0,3)
				for _,v := range t.Players{
					tmp = append(tmp,v)
				}

				index := rand.Intn(len(tmp))
				p := tmp[index]
				p.Info.Gold += t.Gold
				t.Gold = 0
			}
		}
	}
}

func RandomGenGold(g *Game){
	n := MapSize-1
	for i:=0;i<n;i++{
		r := rand.Intn(MapSize)

		x := rand.Intn(MapSize)
		y := rand.Intn(MapSize)

		t := g.Tilemap[y][x]
		t.Gold += r
	}
}

func PlayGameRounds(game *Game){
	initGame(game)
	for {
		if len(prepares) == 0 {
			break
		}

		//random Gold
		RandomGenGold(game)

		//play loop
		pubGameMap(game)

		time.Sleep(time.Second/4)

		//wait for move
		for token, v := range connections {
			msg, ok := <-v.rc
			if ok && msg.Msgtype == 4 {

				log.Println("player is moving:", token)
				LogStruct(msg)

				MovePlayer(game,v,msg.X,msg.Y)
			} else {
				KickPlayer(v)
			}
		}

		game.RoundID++
		ApplyGameLogic(game)

		if CheckGameOver(game){
			game.status = -1
			log.Println("Game Over.")
			SaveGameResult(game)
			break;
		}
	}
}

func GameLoop() {
	game := &Game{status: -1, Msgtype: 3}

	for {
		log.Println("will prepare for next game.")
		for {
			if len(prepares) <= 0 {
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

func SaveGameResult(g *Game){
	ps := make([]*Player,0,len(connections))

	for _, v := range connections {
		ps = append(ps,v)
	}

	sort.Slice(ps, func(i, j int) bool { return ps[i].Info.Gold > ps[j].Info.Gold })

	if len(ps)>0{
		//firt place get 5 score
		p := ps[0]
		records.Scores[p.Info.Key] = records.Scores[p.Info.Key] + 5
	}

	if len(ps)>1{
		//sec place get 2 score
		p := ps[1]
		records.Scores[p.Info.Key] = records.Scores[p.Info.Key] + 2
	}

	if len(ps)>2{
		//sec place get 2 score
		p := ps[2]
		records.Scores[p.Info.Key] = records.Scores[p.Info.Key] + 1
	}

	LogStruct(ps)
	LogStruct(records)
}

type GameRank struct{
	// Gameresults []

	Scores map[string]int
}

func home(w http.ResponseWriter, r *http.Request) {

	out := ""

	bt, err := json.Marshal(records)
	if err!=nil{
		out = "error"
	}else{
		out = string(bt)
	}

	_, _ = w.Write([]byte(out))
}

var records *GameRank

func main() {
	records = &GameRank{}
	records.Scores = make(map[string]int)

	go GameLoop()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)
	// http.HandleFunc("/", home)
	log.Println("ws server ready...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
}
