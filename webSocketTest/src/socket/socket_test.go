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
	"time"
	"log"
	"fmt"
	"testing"
)



func TestSimple(t *testing.T){

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

	var done = make(chan bool)

	go func(){

		var connectTimes int = 0
		var err = fmt.Errorf("no connect")
		var ws *websocket.Conn = nil
		for err != nil {
			connectTimes++
			origin := "http://localhost/"
			url := "ws://localhost:7777/srv"
			ws, err = websocket.Dial(url, "", origin)
			if connectTimes > 20 {
				t.Errorf("have try to connect 20 times")
				done <- true
				return
			}
			if err != nil {
				log.Fatal(err)
			}
		}

		for {

			time.Sleep(1*time.Second)
			var msgToSend  = []byte("hello, world!\n")
			nWrite, err := ws.Write(msgToSend)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("send %d data\n", nWrite)

			var msgReceived = make([]byte, 512)
			var nRead int
			if nRead, err = ws.Read(msgReceived); err != nil {
				log.Fatal(err)
			}
			checkByteSliceEqual(msgToSend, msgReceived[:nRead], t)

			done <- true

		}
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

