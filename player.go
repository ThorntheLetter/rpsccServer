package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var DisconnectChannel = make(chan *player, 5)

type opponent struct {
	scorevs int
	played  bool
	p       *player
}

type player struct {
	name       string
	busy       bool
	connection net.Conn
	reader     *bufio.Reader
	opponents  []opponent
}

func NewOpponent(p player) *opponent {
	o := new(opponent)
	o.p = &p
	o.played = false
	return o
}

func NewPlayer(c net.Conn, opponent []player) *player {
	p := new(player)
	p.busy = true
	p.connection = c
	p.reader = bufio.NewReader(c)
	var err error
	p.name, err = p.reader.ReadString('\n')
	p.name = p.name[:len(p.name)-1]
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, o := range opponent {
		p.AddOpponent(o)
	}
	return p
}

func (p *player) Score() int {
	sum := 0
	for _, o := range p.opponents { // this range shold actually work like that
		sum = sum + o.scorevs
	}
	return sum
}

func (p *player) AddOpponent(newOpp player) {
	o := NewOpponent(newOpp)
	p.opponents = append(p.opponents, *o)
}

func (p *player) PlayLoop() {
	source := rand.NewSource(time.Now().UnixNano())
	randoms := rand.New(source)
	var start time.Time
	var dur time.Duration
	for {
		start = time.Now()
		p.SearchMatch()
		dur = time.Since(start)
		time.Sleep(dur + time.Second + time.Second*time.Duration(randoms.Int63n(2))) //sleep for 0-3 times as long as it took so it can be challenged
		//fmt.Printf("%s: %t\n", p.name, p.busy)
	}
}

func (p *player) SearchMatch() {
	p.busy = true
	names := ""
	for _, s := range p.opponents {
		names = names + " " + s.p.name
	}
	for i, _ := range p.opponents {
		fmt.Printf(p.name + " : " + names + " | " + p.opponents[i].p.name)
		if !p.opponents[i].played {
			fmt.Printf(" is not played")
			if p.opponents[i].p.busy {
				fmt.Printf(", not busy")
				p.Challenge(p.opponents[i])
				p.opponents[i].played = true
				break
			}
		}
		fmt.Println()
	}
	p.busy = false
}

func (p1 *player) Challenge(p1op opponent) {
	fmt.Println(p1.name + " vs " + p1op.p.name)
	p2 := p1op.p
	p2.busy = true

	var p2op opponent

	for i, _ := range p2.opponents {
		if p2.opponents[i].p == p1 {
			p2op = p2.opponents[i]
			break
		}
	}

	p1op.played = true
	p2op.played = true

	p1.connection.Write([]byte("6\n"))
	p2.connection.Write([]byte("6\n"))

	for i := 0; i < 99; i += 1 {
		p1in, err := p1.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}
		p2in, err := p2.reader.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
		}

		Winner(p1op, p1in, p2op, p2in)

		_, err = p1.connection.Write([]byte(p2in))
		if err != nil {
			fmt.Println(err.Error())
		}
		_, err = p2.connection.Write([]byte(p1in))
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	p2.busy = false
}

func Winner(p1op opponent, p1option string, p2op opponent, p2option string) {
	switch p1option {
	case p2option:
		return
	case "1\n":
		switch p2option {
		case "3\n", "4\n", "5\n":
			p1op.scorevs += 1
			return
		case "2\n":
			p2op.scorevs += 1
			return
		}
	case "2\n":
		switch p2option {
		case "1\n":
			p1op.scorevs += 1
			return
		case "3\n", "4\n", "5\n":
			p2op.scorevs += 1
			return
		}
	case "3\n":
		switch p2option {
		case "2\n", "5\n":
			p1op.scorevs += 1
			return
		case "1\n", "4\n":
			p2op.scorevs += 1
			return
		}
	case "4\n":
		switch p2option {
		case "2\n", "3\n":
			p1op.scorevs += 1
			return
		case "1\n", "5\n":
			p2op.scorevs += 1
			return
		}
	case "5\n":
		switch p2option {
		case "2\n", "4 \n":
			p1op.scorevs += 1
			return
		case "1\n", "3\n":
			p2op.scorevs += 1
			return
		}

	}
}
