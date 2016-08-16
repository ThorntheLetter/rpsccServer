package main

import (
	"fmt"
	"net/http"
)

func root(w http.ResponseWriter, r *http.Request) {
	for i, _ := range playerList {
		w.Write([]byte(fmt.Sprintf("%s:\t%d\n", playerList[i].name, playerList[i].Score())))
	}
}
