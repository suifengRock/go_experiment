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
		fmt.Printf(err.Error())
	}
}

func StartServer(){
	Serve()
}

type ClientId int64
func genClientId() ClientId{
	return 1
}

var clients = make(map[ClientId]*websocket.Conn)

func websocketHandler2(ws *websocket.Conn){

	id := genClientId()
	clients[id] = ws
	if onConnectHandler != nil {
		onConnectHandler(id)
	}

	for{
		//todo: should ensure the buffer is large enough
		var msg []byte = make([]byte, 512)
		nRead, errRead := ws.Read(msg)

		if nRead > 0 {
			onMessageHandler(id, msg[:nRead])
		}

		if errRead != nil {
			fmt.Printf("read error on client id %d\n", id)
			return
		}

		if nRead <= 0 {
			//if read error occur,the nRead value will be zero
			fmt.Printf("has no error on read but read %d data\n", nRead)
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
		}

	}
}







