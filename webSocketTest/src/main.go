/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 10/3/13
 * Time: 4:54 PM
 * To change this template use File | Settings | File Templates.
 */
package main

import (
	"socket"
	"container/list"
)

func echoServer() {

		var clients *list.List = list.New()

		socket.OnConnect(func(clientId socket.ClientId){
			clients.PushBack(clientId)
		})

		socket.OnDisconnect(func(clientId socket.ClientId){

			for e := clients.Front(); e != nil; e = e.Next(){
				if e.Value.(socket.ClientId) == clientId {
					clients.Remove(e)
				}

			}
		})

		socket.OnMessage(func(clientId socket.ClientId, msg []byte){
			socket.Send([]socket.ClientId{clientId}, msg)
		})

		socket.StartServer()



}

func main() {
 	echoServer()
}
