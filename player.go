package main

import (
	"bufio"
	"fmt"
	"net"
)

var DisconnectSlice = []int{}

type opponent struct {
	scorevs int
	p       *player
}

type player struct {
	name       string
	connection net.Conn
	reader     *bufio.Reader
	id         int
	opponents  []*opponent
}

func NewOpponent(p *player) *opponent {
	o := new(opponent)
	o.p = p
	return o
}

func NewPlayer(c net.Conn, id int) *player {
	p := new(player)
	p.connection = c
	p.reader = bufio.NewReader(c)
	p.id = id
	var err error
	p.name, err = p.reader.ReadString('\n')
	p.name = p.name[:len(p.name)-1]
	if err != nil {
		fmt.Println(err.Error())
	}
	return p
}

func (p *player) Score() int {
	sum := 0
	for i, _ := range p.opponents {
		sum += p.opponents[i].scorevs
	}
	return sum
}

func (p *player) AddOpponent(newOpp *player) *opponent {
	o := NewOpponent(newOpp)
	p.opponents = append(p.opponents, o)
	return o
}

func (p1 *player) Challenge(p2 *player) {
	disc := false
	fmt.Println(p1.name + " vs " + p2.name)
	p1op := p1.AddOpponent(p2)
	p2op := p2.AddOpponent(p1)

	p1.connection.Write([]byte("6\n"))
	p2.connection.Write([]byte("6\n"))

	for i := 0; i < 100; i += 1 {
		p1in, err := p1.reader.ReadString('\n')
		if err != nil {
			fmt.Println(p1.name, "disconnected")
			DisconnectSlice = append(DisconnectSlice, p1.id)
			disc = true
			p1in = "0\n"

		}
		p2in, err := p2.reader.ReadString('\n')
		if err != nil {
			fmt.Println(p2.name, "disconnected")
			DisconnectSlice = append(DisconnectSlice, p2.id)
			disc = true
			p2in = "0\n"
		}

		if !disc {
			Winner(p1op, p1in, p2op, p2in)
		}

		p1.connection.Write([]byte(p2in)) // any errors here should also get caught and dealt with above,
		p2.connection.Write([]byte(p1in))

		if disc {
			break
		}

	}
}

func disconnect(pList []player, disconnected int) []player {
	for i, _ := range pList {
		for j, _ := range pList[i].opponents {
			if pList[i].opponents[j].p.id == disconnected {
				pList[i].opponents = append(pList[i].opponents[:j], pList[i].opponents[j+1:]...) // remove disconnected player from list of opponents
				break
			}
		}
	}
	for i, _ := range pList {

		if pList[i].id == disconnected {
			pList = append(pList[:i], pList[i+1:]...) // remove disconnected player so it doesnt screw with loop
			break
		}
	}

	return pList
}

func Winner(p1op *opponent, p1option string, p2op *opponent, p2option string) {
	switch p1option {
	case p2option:

	case "1\n":
		switch p2option {
		case "3\n", "4\n", "5\n":
			p1op.scorevs += 1
		case "2\n":
			p2op.scorevs += 1
		}

	case "2\n":
		switch p2option {
		case "1\n":
			p1op.scorevs += 1
		case "3\n", "4\n", "5\n":
			p2op.scorevs += 1
		}

	case "3\n":
		switch p2option {
		case "2\n", "5\n":
			p1op.scorevs += 1
		case "1\n", "4\n":
			p2op.scorevs += 1
		}

	case "4\n":
		switch p2option {
		case "2\n", "3\n":
			p1op.scorevs += 1
		case "1\n", "5\n":
			p2op.scorevs += 1
		}

	case "5\n":
		switch p2option {
		case "2\n", "4\n":
			p1op.scorevs += 1
		case "1\n", "3\n":
			p2op.scorevs += 1
		}
	}
}

type sortablePlayerSlice []player

func (l sortablePlayerSlice) Len() int {
	return len(l)
}

func (l sortablePlayerSlice) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l sortablePlayerSlice) Less(i int, j int) bool {
	return l[i].Score() > l[j].Score() // actually greater than so it sorts descending
}
