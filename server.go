package main

import (
	"fmt"
	"net"
	"net/http"
)

const gamePort = ":7633" // i have not yet seen these used as an integer or without :, so they are strings for now
const webPort = ":7655"

var playerList []player

func main() {
	// network setup
	fmt.Println("server starting")
	ln, err := net.Listen("tcp", gamePort)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ln.Close()

	// website setup
	http.HandleFunc("/", root)
	go http.ListenAndServe(webPort, nil)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		newP := NewPlayer(conn)
		fmt.Println("connection accepted: " + newP.name)
		for i, _ := range playerList {
			newP.Challenge(&playerList[i])
		}
		playerList = append(playerList, *newP)

	}
}
