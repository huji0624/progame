package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func WhateverOrigin(r *http.Request) bool {
	return true
}

/*
	msgtype 0登陆 -1登陆失败 1请准备 2准备好了 3游戏信息 4玩家行动
*/
type Msg struct {
	Msgtype int
	Token   string
	X       int
	Y       int
}

type PlayerInfo struct {
	X    int
	Y    int
	Gold int
}

type Player struct {
	c  *websocket.Conn
	rc chan *Msg

	Info *PlayerInfo
}

var upgrader = websocket.Upgrader{CheckOrigin: WhateverOrigin} // use default options
var connections = make(map[string]*Player)
var connections_count int

func LogStruct(v interface{}) {
	bt, err := json.Marshal(v)
	if err == nil {
		log.Println(string(bt))
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
		log.Println("write data:")
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

	outtime := time.Now().Add(1 * time.Second)
	c.SetReadDeadline(outtime)

	jmsg := Msg{}
	err = WaitForMsg(c, &jmsg)
	if err != nil {
		log.Println("read err:", err)
	} else {
		if jmsg.Token == "jaisjis" {
			rc := make(chan *Msg)
			p := &Player{c: c, rc: rc}
			connections[jmsg.Token] = p
			connections_count++
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
		outtime := time.Now().Add(1 * time.Second)
		p.c.SetReadDeadline(outtime)

		jmsg := Msg{}
		err := p.c.ReadJSON(&jmsg)
		if err == nil {
			p.rc <- &jmsg
		} else {
			log.Println("read token:", token, " err:", err)
			KickPlayer(token)
			break
		}
	}
}

func KickPlayer(token string) {
	p := connections[token]
	p.c.Close()
	delete(connections, token)
	connections_count--
	log.Println("kick player:", token)
}

func home(w http.ResponseWriter, r *http.Request) {

}

type Tile struct {
	Gold    int
	Players []*Player
}

const (
	MapSize = 3
)

type Game struct {
	Msgtype int
	status  int //-1无效0准备1开始

	Tilemap [MapSize][MapSize]*Tile
}

func initGame(g *Game) {
	log.Println("init game.")
	size := MapSize
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			t := &Tile{Gold: 0}
			t.Players = make([]*Player, 0)
			g.Tilemap[i][j] = t
		}
	}

	for _, v := range connections {
		x := rand.Intn(MapSize)
		y := rand.Intn(MapSize)
		v.Info = &PlayerInfo{X: x, Y: y}
		v.Info.Gold = 0

		t := g.Tilemap[y][x]
		t.Players = append(t.Players, v)
	}
}

func PublicToClientData(data []byte) {
	log.Println("will pub data:", string(data))
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

func GameLoop() {
	game := &Game{status: -1, Msgtype: 3}

	for {
		log.Println("will prepare for next game.")
		for {
			if connections_count <= 0 {
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
		for _, v := range connections {
			WriteToClient(v.c, jmsg)
		}

		time.Sleep(time.Second)

		log.Println("wait for response.")
		//prepare response
		jmsg = Msg{}
		for token, v := range connections {
			msg, ok := <-v.rc
			if ok && msg.Msgtype == 2 {
				log.Println("player is ready:", token)
			} else {
				KickPlayer(token)
			}
		}

		if connections_count == 0 {
			continue
		}

		game.status = 1
		initGame(game)
		for {
			//play loop
			pubGameMap(game)

			time.Sleep(time.Second)

			//wait for move
			for token, v := range connections {
				msg, ok := <-v.rc
				if ok && msg.Msgtype == 4 {
					log.Println("player is moving:", token)
					LogStruct(msg)
				} else {
					KickPlayer(token)
				}
			}

		}
	}
}

func main() {
	go GameLoop()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	// http.HandleFunc("/", home)
	log.Println("ws server ready...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
}
