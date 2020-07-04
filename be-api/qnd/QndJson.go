package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type MyStruct struct {
		Id    string `json:"id"`
		Value int    `json:"value"`
	}
	js := []byte(`{"id":"1","value":23}`)
	s := &MyStruct{}
	err := json.Unmarshal(js, &s)
	fmt.Println(err)
	fmt.Println(s)
}
