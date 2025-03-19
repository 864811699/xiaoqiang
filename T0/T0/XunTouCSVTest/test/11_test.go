package test

import (
	"fmt"
	"reflect"
	"testing"
)

type Request struct {
	basic
	Stockcode           string
}

type basic struct {
	Sym, Qid string
}

func StructToMap(obj interface{}) map[string]string {

	objVal := reflect.ValueOf(obj)
	if objVal.Kind()!=reflect.Struct{
		fmt.Println("非结构体")
		return nil
	}
	objType := reflect.TypeOf(obj)
	var data = make(map[string]string)

	for i:=0;i<objType.NumField();i++{
		fieldType:=objType.Field(i)
		fieldVal:=objVal.Field(i)
		if fieldVal.Kind()==reflect.Struct&&fieldVal.CanInterface(){
			m:=StructToMap(fieldVal.Interface())
			for k,v:=range m{
				data[k]=v
			}
		}else {
			data[fieldType.Name]=fmt.Sprintf("%v",fieldVal.Interface())
		}

	}

	//for i := 0; i < objType.NumField(); i++ {
	//	fmt.Println(objVal)
	//	if objVal.Field(i).CanInterface() {
	//		data[objType.Field(i).Name] = fmt.Sprintf("%v", objVal.Field(i).Interface())
	//	}
	//}
	return data
}

func TestRe(t *testing.T)  {
	r:=Request{
		basic:     basic{Sym:"aaa",Qid:"bbb"},
		Stockcode: "ccc",

	}
	m:=StructToMap(r)
	t.Log(m)
}