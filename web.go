package main

import (
	"fmt"
	"net/http"
	"sort"
)

func root(w http.ResponseWriter, r *http.Request) {
	var webPlayers []player
	copy(webPlayers, playerList) // copy so sorting doesnt effect the playing
	sort.Sort(sortablePlayerSlice(playerList))

	w.Write([]byte("<table>"))
	for i, _ := range playerList {
		w.Write([]byte(fmt.Sprintf("<tr><th>%s</th><th>%d<tr>", playerList[i].name, playerList[i].Score())))
	}
}
