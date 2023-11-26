package test

import (
	"fmt"
	"github.com/aivyss/typex/pointer"
	"reflect"
	"testing"
)

func TestTemp(t *testing.T) {
	var str *string = pointer.MustPointer("aaa")
	fmt.Println("*string reflect.TypeOf(str)=", reflect.TypeOf(str))
	fmt.Println("*string reflect.TypeOf(str).Kind()=", reflect.TypeOf(str).Kind())
	fmt.Println("*string reflect.TypeOf(str).Kind() == reflect.Pointer", reflect.TypeOf(str).Kind() == reflect.Pointer)
	fmt.Println("*string reflect.TypeOf(str).Elem()=", reflect.TypeOf(str).Elem().Kind() == reflect.String)
}
