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
)


func startEchoServer(){
	go func(){

		var clients *list.List = list.New()

		OnConnect(func(clientId ClientId){
			clients.PushBack(clientId)
		})

		OnDisconnect(func(clientId ClientId){

			for e := clients.Front(); e != nil; e = e.Next(){
				if e.Value.(ClientId) == clientId {
					clients.Remove(e)
				}

			}
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

func TestSimple(t *testing.T){

	startEchoServer()

	var done = make(chan bool)

	go func(){

		ws, err := NewClient()

		if err != nil {
			log.Fatal(err)
			done<-true
			return
		}

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

		done<-true
	}()

	<-done

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

