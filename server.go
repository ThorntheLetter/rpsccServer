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

		newP := NewPlayer(conn)
		fmt.Println("connection accepted: " + newP.name)
		for i, _ := range playerList {
			newP.Challenge(playerList[i])
		}
		playerList = append(playerList, *newP)
		// go newP.PlayLoop()

	}
}
