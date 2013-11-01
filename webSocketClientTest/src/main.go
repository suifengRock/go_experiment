/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 10/3/13
 * Time: 8:42 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"fmt"
	"time"
	"math/rand"
	"log"
	"code.google.com/p/go.net/websocket"
)

var val int = 0

func test()int{
	val++
	fmt.Println("val %d\n", val)
	return val
}

func testClient(){
	origin := "http://localhost/"
	url := "ws://localhost:7777/srv"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(5*time.Second)
		nWrite, err := ws.Write([]byte("hello, world!\n"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("send %d data\n", nWrite)
//		var msg = make([]byte, 512)
//		var n int
//		if n, err = ws.Read(msg); err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(msg[:n])
//		fmt.Printf("Received: %s.\n", msg[:n])
	}

}

func main() {
	testClient()
}

func randomDuration() time.Duration {
	return time.Duration(rand.Int63n(1e9))
}

