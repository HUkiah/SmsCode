package main

import (
	"log"
	"reflect"
)

type T struct {
	A int
	B string
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func main() {
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem()
	// typeOfT := s.Type()
	value := make([]interface{}, 0)
	value[0] = s.FieldByName().Addr().Interface()
	Log(value)

	// for i := 0; i < s.NumField(); i++ {
	// 	f := s.Field(i)
	// 	fmt.Printf("%d: %s %s = %v\n", i,
	// 		typeOfT.Field(i).Name, f.Type(), f.Interface())
	// }
	// s.Field(0).SetInt(77)
	// s.Field(1).SetString("Sunset Strip")
	// fmt.Println("t is now", t)
}
