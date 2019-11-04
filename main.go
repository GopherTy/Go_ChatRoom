package main

import (
	"chatroom/cmd"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang/protobuf/proto"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
