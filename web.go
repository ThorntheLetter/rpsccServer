package main

import (
	"fmt"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<table>"))
	for i, _ := range playerList {
		w.Write([]byte(fmt.Sprintf("<tr><th>%s</th><th>%d<tr>", playerList[i].name, playerList[i].Score())))
	}
}
