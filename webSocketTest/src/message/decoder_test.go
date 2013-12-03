/**
 * Created with IntelliJ IDEA.
 * User: Administrator
 * Date: 12/1/13
 * Time: 2:24 PM
 * To change this template use File | Settings | File Templates.
 */
package message

import (
	"testing"
	"log"
)

type Input struct {
	A *int
	B *int
}

func (input *Input) GetA() int {
	return *input.A
}

func (input *Input) GetB() int {
	return *input.B
}

func Invoke(a int, b int) int {
	return a+b
}

type Output struct {
	A *int
}

func TestAdapterInvoke(t *testing.T) {
	A := 1
	B := 2
	inputObj := Input{&A, &B}
 	v := AdapterInvoke(&inputObj, Invoke)
	if len(v) != 1 {
		t.Errorf("return param count should be 1 got:%v", len(v))
	}
	r := v[0].(int)
	if r != 3 {
		t.Errorf("return value should be 3 got:%v", r)
	}
}

func TestAdapterInvokeReturn(t *testing.T){
	log.Printf("start")
	A := 1
	B := 2
	inputObj := Input{&A, &B}
	outputObj := Output{}
	AdapterInvokeReturn(&inputObj, &outputObj, Invoke)
	if *outputObj.A != 3 {
		t.Errorf("output obj.A should be 3 but got:%v", *outputObj.A)
	}
}
