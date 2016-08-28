package main

import (
	"fmt"
	"net"
	"net/http"
	"sort"
)

func root(w http.ResponseWriter, r *http.Request) {
	webPlayers := make(sortablePlayerSlice, len(playerList))
	fmt.Println(copy(webPlayers, playerList)) // copy so sorting doesnt effect the playing
	sort.Sort(sortablePlayerSlice(webPlayers))

	w.Write([]byte("<table>"))
	for _, p := range webPlayers {
		w.Write([]byte(fmt.Sprintf("<tr><th>%s</th><th>%d<tr>", p.name, p.Score())))
	}
}

func cody() {
	ln, err := net.Listen("tcp", ":3489")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}

		webPlayers := make(sortablePlayerSlice, len(playerList))
		copy(webPlayers, playerList) // copy so sorting doesnt effect the playing
		sort.Sort(webPlayers)
		for _, p := range webPlayers {
			// fmt.Printf("%s: %d", p.name, p.Score())
			conn.Write([]byte(fmt.Sprintf("%s: %d\n", p.name, p.Score())))
		}
		conn.Close()
	}
}
