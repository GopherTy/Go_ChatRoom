package main

import (
	"log"
	"chatroom/cmd"

	_ "github.com/golang/protobuf/proto"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
