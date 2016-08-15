package main

import (
	"fmt"
	"net"
)

func main() {
	// network setup
	fmt.Println("server starting")
	ln, err := net.Listen("tcp", ":7633")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ln.Close()

	var playerList []player

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		newP := NewPlayer(conn, playerList)
		fmt.Println("connection accepted: " + newP.name)
		for i, _ := range playerList {
			playerList[i].AddOpponent(*newP)
		}
		playerList = append(playerList, *newP)
		go newP.PlayLoop()

	}
}
