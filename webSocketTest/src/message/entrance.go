/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 11/13/13
 * Time: 10:29 PM
 * To change this template use File | Settings | File Templates.
 */
package message

import (
	"logic"
)

var actionMap = make(map[string]func(interface{})[]byte)
var actionEncodeMap =  make(map[string]func([]byte)interface{})

func Init(){
	actionMap["Move"] = actionMoveEncode
	actionEncodeMap["Move"] = actionMoveDecode
}

func DecodePackage(msg []byte) (string, interface{}){
	fn := actionEncodeMap["test"]
	return nil, fn(msg)
}

func EncodeMessage(action string, obj interface{}) []byte{
	//lookup obj via action
	fn := actionMap[action]
	return fn(obj)

}

func actionMoveEncode(obj interface{})[]byte{
	role := obj.(logic.Role)
	x := role.X
	y := role.Y
	return nil
}

func actionMoveDecode(msg []byte)interface {}{
	return logic.Role{}
}



