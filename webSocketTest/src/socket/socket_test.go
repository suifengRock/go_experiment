/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 10/6/13
 * Time: 10:54 PM
 * To change this template use File | Settings | File Templates.
 */
package socket

import (
	"container/list"
	"code.google.com/p/go.net/websocket"
	"log"
	"fmt"
	"testing"
	"time"
	"sync"
)

var once sync.Once

type ClientList struct {
	mutex sync.Mutex
	data *list.List
}

func NewClientList() *ClientList {
	return &ClientList{data:list.New()}
}

func (clientList *ClientList) Add(id ClientId) error{
	clientList.mutex.Lock()
	clientList.data.PushBack(id)
	clientList.mutex.Unlock()
	return nil
}

func (clientList *ClientList) Remove(id ClientId) error {
	clientList.mutex.Lock()
	for e := clientList.data.Front(); e != nil; e = e.Next(){
		if e.Value.(ClientId) == id {
			clientList.data.Remove(e)
		}
	}
	clientList.mutex.Unlock()
	return nil
}

func (clientList *ClientList) Length() int {
	return clientList.data.Len()
}

var testClients = NewClientList()
func startEchoServer(){

	go func(){

		OnConnect(func(clientId ClientId){
			testClients.Add(clientId)
		})

		OnDisconnect(func(clientId ClientId){
			log.Println("dis conncet")
			testClients.Remove(clientId)
		})

		OnMessage(func(clientId ClientId, msg []byte){
			Send([]ClientId{clientId}, msg)
		})

		StartServer()


	}()

}

type TestClient struct {
	 *websocket.Conn
}

func NewClient() (*TestClient,error) {

	var connectTimes int = 0
	var err = fmt.Errorf("no connect")
	var ws *websocket.Conn = nil
	for {
		connectTimes++
		origin := "http://localhost/"
		url := "ws://localhost:7777/srv"
		ws, err = websocket.Dial(url, "", origin)

		if ws != nil {
			break
		}

		if connectTimes > 20 {
			return nil, fmt.Errorf("hava try connect 20 times")
		}
		if err != nil {
			fmt.Printf(err.Error())
			time.Sleep(time.Second)
		}
	}
	return &TestClient{ws},nil
}

func (ws *TestClient) Read(msg []byte) (n int, err error) {
	return ws.Conn.Read(msg)
}

func (ws *TestClient) Write(msg []byte) (n int, err error) {
	return ws.Conn.Write(msg)
}

func clientConnectToServer(connected chan bool, checkConnected chan bool, done chan bool, t *testing.T){
		ws, err := NewClient()

		if err != nil {
			log.Fatal(err)
			done<-true
			return
		}

		connected <- true
		<-checkConnected

		var msgToSend  = []byte("hello, world!\n")
		nWrite, err := ws.Write(msgToSend)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("send %d data\n", nWrite)

		var msgReceived = make([]byte, 512)
		var nRead int
		if nRead, err = ws.Read(msgReceived); err != nil {
			log.Fatal(err)
		}
		checkByteSliceEqual(msgToSend, msgReceived[:nRead], t)
		ws.Conn.Close()
		done<-true

}

func TestSimple(t *testing.T){

	once.Do(startEchoServer)

	var done = make(chan bool)
	var connected = make(chan bool)
	var checkConnected = make(chan bool)

	var countClients = 10

	for i:=0; i<countClients; i++ {
		go clientConnectToServer(connected, checkConnected, done, t)
	}

	for i:=0; i<countClients; i++ {
		<-connected
	}

	if testClients.Length() != countClients {
		t.Errorf("should hava %v connected testClients, got:%v", countClients, testClients.Length())
	}

	for i:=0; i<countClients; i++ {
		checkConnected <- true
	}

	for i:=0; i<countClients; i++ {
		<-done
	}

	if countClients != countClients-testClients.Length() {
		t.Errorf("should have closed %v connected testClient, got:%v", countClients, countClients-testClients.Length())
	}

}

func checkByteSliceEqual(expect []byte, got []byte, t *testing.T){
	if len(expect) != len(got) {
		t.Errorf("msg len is not equal expect:%v, got:%v", len(expect), len(got))
	}
	for i, v := range expect {
		if v != got[i] {
			t.Errorf("msg not equal expect:%v, got:%v", expect, got)
		}
	}

}

