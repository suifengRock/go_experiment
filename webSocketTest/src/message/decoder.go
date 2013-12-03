/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 12/1/13
 * Time: 2:23 PM
 * To change this template use File | Settings | File Templates.
 */
package message

import (
	"reflect"
)

func AdapterInvoke(obj interface {}, fn interface {})[]interface {}{

	vPtr := reflect.ValueOf(obj)
	v := vPtr.Elem()
	out := make([]reflect.Value, v.NumField())
	for i:=0; i<v.NumField(); i++ {
		method := vPtr.MethodByName("Get"+v.Type().Field(i).Name)
		out[i] = method.Call([]reflect.Value{})[0]
	}
	fnVal := reflect.ValueOf(fn)
	rVal :=  fnVal.Call(out)

	interfaceVal := make([]interface {}, fnVal.Type().NumOut())
	for i:=0; i<fnVal.Type().NumOut(); i++ {
		interfaceVal[i] = rVal[i].Interface()
	}
	return interfaceVal

}

func AdapterInvokeReturn(in interface {}, out interface {}, fn interface {}){

//	v := reflect.ValueOf(in)
//	outVal := make([]reflect.Value, v.NumField())
//	for i:=0; i<v.NumField(); i++ {
//		method := v.MethodByName("Get"+v.Type().Field(i).Name)
//		outVal[i] = method.Call([]reflect.Value{})[0]
//	}
//	fnVal := reflect.ValueOf(fn)
//	rVal :=  fnVal.Call(outVal)

	inInterface := AdapterInvoke(in, fn)
	valOutPtr := reflect.ValueOf(out)
	valOut := valOutPtr.Elem()
	valFn := reflect.ValueOf(fn)
	for i:=0; i<valOut.NumField(); i++ {
		valOut.Field(i).Set(reflect.ValueOf(&inInterface[0]))
	}

}

