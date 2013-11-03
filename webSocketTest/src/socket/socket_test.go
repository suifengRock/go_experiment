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
			log.Println("disconncet")
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

var testMutex sync.Mutex
func TestSimple(t *testing.T){

	testMutex.Lock()
	defer testMutex.Unlock()

	once.Do(startEchoServer)

	var done = make(chan bool)
	var connected = make(chan bool)
	var checkConnected = make(chan bool)

	var countClients = 10

	for i:=0; i<countClients; i++ {
		go clientConnectToServer(connected, checkConnected, done, t)
	}

	//wait all client connected to server
	for i:=0; i<countClients; i++ {
		<-connected
	}

	if testClients.Length() != countClients {
		t.Errorf("should hava %v connected testClients, got:%v", countClients, testClients.Length())
	}

	//inform all client to continue
	for i:=0; i<countClients; i++ {
		checkConnected <- true
	}

	//wait all client to test complete
	for i:=0; i<countClients; i++ {
		<-done
	}

	//have enough time for server to close connection
	time.Sleep(5*time.Second)
	//check should all client have closed connection
	if countClients != countClients-testClients.Length() {
		t.Errorf("should have closed %v connected testClient, got:%v", countClients, countClients-testClients.Length())
	}

}

func TestHeartBreak(t *testing.T){

	testMutex.Lock()
	defer testMutex.Unlock()

	var heartbreakData = []byte("a")

	once.Do(startEchoServer)

	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	if testClients.Length() != 1 {
		t.Errorf("should hava %v clients but got %v", 1, testClients.Length())
	}

	//send a byte length msg indicate that this is a heart break package
	for i:=0; i<5; i++ {
		_, err := client.Write(heartbreakData)
		if err != nil {
			t.Errorf("should send heartbreak data sucess")
		}
		time.Sleep(2*time.Second)
	}

	if testClients.Length() != 1 {
		t.Errorf("should hava %v clients but got %v", 1, testClients.Length())
	}

	time.Sleep(10*time.Second)
	log.Printf("try to send msg to server after heart break timeout \n")

	if testClients.Length() != 0 {
		t.Errorf("should hava %v clients but got %v", 0, testClients.Length())
	}

}




