/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 10/5/13
 * Time: 10:06 PM
 * To change this template use File | Settings | File Templates.
 */
package socket

import (
	"net/http"
	"code.google.com/p/go.net/websocket"
	"time"
	"fmt"
	"math"
	"log"
	"uuid"
)

var clientConnect chan bool = make(chan bool)
var clientDisConnect chan bool = make(chan bool)

func websocketHandler(ws *websocket.Conn){

	defer func(){
		ws.Close()
	}()
	clientConnect <- true
	read := make(chan bool)
	go heartBreak(read)

	for {
		time.Sleep(time.Second)
		var msg []byte = make([]byte, 512)
		nRead, errRead := ws.Read(msg)
		if nRead > 0 {
			read <- true
		}
		if errRead != nil {
			return
		}
	}
}

func heartBreak(read chan bool){
	for {
		select{
		case <-read:
		case <-time.After(5*time.Second):
			clientDisConnect <- true
			return
		}
	}
}

func benchmark(connect chan bool, disconnect chan bool){

	//时间 最大连接数 最小连接数
	var maxClients int64 = 0
	var minClients int64 = math.MaxInt64
	var total int64 = 0

	for{
		select {
		case <-connect:
			total++
			if total > maxClients {
				maxClients = total
			}
		case <-disconnect:
			total--
			if total < minClients {
				minClients = total
			}
		case <-time.After(1 * time.Minute):
			fmt.Printf("run %d minutes,total:%d, max:%d, min:%d\n", 5, total, maxClients, minClients)
			return
		}
	}


}

func Serve(){

	http.Handle("/srv", websocket.Handler(websocketHandler2))
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal(err)
	}
}

func StartServer(){
	Serve()
}

type ClientId *uuid.UUID
func genClientId() (ClientId, error){
	u4, err := uuid.NewV4()
	if err != nil {
		return nil,err
	}
	return u4,nil
}

var clients = make(map[ClientId]*websocket.Conn)

func heartBreak2(read chan bool, timeout chan bool){
	for {
		select{
		case <-read:
		case <-time.After(3*time.Second):
			timeout <- true
			return
		}
	}
}

func websocketHandler2(ws *websocket.Conn){

	id, err := genClientId()

	if err != nil {
		log.Fatal(err)
	}

	clients[id] = ws
	if onConnectHandler != nil {
		onConnectHandler(id)
	}

	defer func(){
		err:=ws.Close()
		if err != nil {
			log.Printf("close ws fail:", err.Error())
		}
		log.Printf("close ws on id %d \n", id)
	}()

	var dataChan = make(chan []byte)
	var heartbreakChan = make(chan bool)
	var onReadErr = make(chan bool)

	go func(){
		for{
			var msg []byte
			errRead := websocket.Message.Receive(ws, &msg)

			if errRead != nil {
				log.Printf("read error on client id %d\n", id)
				onReadErr <- true
				return
			}

			if len(msg) == 1 {
				heartbreakChan <- true
			}
			dataChan<-msg

		}
	}()

	for{
		select {
		case msg := <-dataChan:
			onMessageHandler(id, msg)
		case <-heartbreakChan:
			log.Println("heart break data recived")
		case <-time.After(3*time.Second):
			log.Println("heart break timeout")
			onDisconnectHandler(id)
			return
		case <-onReadErr:
			onDisconnectHandler(id)
			log.Printf("on read err chan got msg")
			return
		}
	}
}

var onConnectHandler func(clientId ClientId)
func OnConnect(handler func(clientId ClientId)){
	onConnectHandler = handler
}

var onMessageHandler func(clientId ClientId, msg []byte)
func OnMessage(handler func(clientId ClientId, msg []byte)){
	onMessageHandler = handler
}

var onDisconnectHandler func(clientId ClientId)
func OnDisconnect(handler func(clientId ClientId)){
	onDisconnectHandler = handler
}

func Send(clientIds []ClientId, msg []byte){
	//todo: what it mean when write err
	for _, clientId := range clientIds {
		if _, exist := clients[clientId]; !exist {
			fmt.Println("client id %d is not exist in client map\n", clientId)
			continue
		}
		ws := clients[clientId]
		nWrite, errWrite := ws.Write(msg)

		if nWrite <= 0 {
			fmt.Printf("warning:write %d data size!!\n", nWrite)
		}
		if errWrite != nil {
			fmt.Printf("client id %d write msg error!!\n", clientId)
			onDisconnectHandler(clientId)
		}

	}
}

type Server struct {
	port string
	clients map[ClientId]*websocket.Conn
	onConnectHandler func(ClientId)
	onMessageHandler func(ClientId,  []byte)
	onDisconnectHandler func(ClientId)
}

func NewServer(port string) *Server {
	var clients = make(map[ClientId]*websocket.Conn)
	return &Server{clients:clients,port:port}
}

func (server *Server)OnConnectHandler(handler func(clientId ClientId)){
	server.onConnectHandler = handler
}

func (server *Server)OnMessageHandler(handler func(clientId ClientId, msg []byte)){
	server.onMessageHandler = handler
}

func (server *Server)OnDisconnectHandler(handler func(clientId ClientId)){
	server.onDisconnectHandler = handler
}

func (server *Server) Send(clientIds []ClientId, msg []byte){
	for _, clientId := range clientIds {
		if _, exist := server.clients[clientId]; !exist {
			fmt.Println("client id %d is not exist in client map\n", clientId)
			continue
		}
		ws := server.clients[clientId]
		nWrite, errWrite := ws.Write(msg)

		if nWrite <= 0 {
			fmt.Printf("warning:write %d data size!!\n", nWrite)
		}
		if errWrite != nil {
			fmt.Printf("client id %d write msg error!!\n", clientId)
		}

	}
}

func (server *Server)websocketHandler(ws *websocket.Conn){
	id, err := genClientId()

	if err != nil {
		log.Fatal(err)
	}

	server.clients[id] = ws
	if server.onConnectHandler != nil {
		server.onConnectHandler(id)
	}

	defer ws.Close()

	for{
		var msg []byte
		errRead := websocket.Message.Receive(ws, &msg)

		if errRead != nil {
			log.Printf("read error on client id %d\n", id)
			server.onDisconnectHandler(id)
			break
		}

		server.onMessageHandler(id, msg)

	}
}

func (server *Server)Start(){
	http.Handle("/srv", websocket.Handler(server.websocketHandler))
	if err := http.ListenAndServe(":"+server.port, nil); err != nil {
		log.Fatal(err)
	}
}









