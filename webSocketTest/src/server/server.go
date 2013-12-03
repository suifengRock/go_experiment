/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 11/14/13
 * Time: 8:07 PM
 * To change this template use File | Settings | File Templates.
 */
package server

import (
	"socket"
	"container/list"
	"log"
	"message"
	"logic"
	"sync"
)

var clients *list.List = list.New()

func Init(){

	socket.OnConnect(func(clientId socket.ClientId){
		log.Printf("client connect")
		clients.PushBack(clientId)
	})

	socket.OnDisconnect(func(clientId socket.ClientId){

		log.Printf("client disconnet")

		for e := clients.Front(); e != nil; e = e.Next(){
			if e.Value.(socket.ClientId) == clientId {
				clients.Remove(e)
			}

		}
	})

	socket.OnMessage(func(clientId socket.ClientId, msg []byte){
		process(msg)
	})

	socket.StartServer()

}

func process(msg []byte){

	action, obj := message.DecodePackage(msg)
	resultOjb := logic.Execute(action, obj)
	socketClientIds := sync.CalRange(action, resultOjb)
	msgToSend := message.EncodeMessage(action, resultOjb)
	socket.Send(socketClientIds, msgToSend)


	"auth.Login"
	LoginRequest object

	mainInterFaceParse := &protobuf.MainInterFace{}
	err := proto.Unmarshal(arg_recMsg, mainInterFaceParse)
	action := mainInterFaceParse.GetActionName()
	v := reflect.ValueOf(mainInterFaceParse)
    requestObj := v.MethodByName("Get"+"Login"+"Request").Call()[0].Interface()

	t :=requestObj.Type()
	val := make([]interface{}, t.NumFiled())
	for i:= 0; i<t.NumFiled(); i++ {
    	filedName := t.FieldByIndex(i).Name()
		val[i] = t.MethodByName("Get"+fieldName).Call()[0].Interface()
	}

}




