package main

import (
	"bufio"
	"fmt"
	"net"
)

var DisconnectChannel = make(chan *player, 5)

type opponent struct {
	scorevs int
	p       *player
}

type player struct {
	name       string
	connection net.Conn
	reader     *bufio.Reader
	opponents  []opponent
}

func NewOpponent(p player) *opponent {
	o := new(opponent)
	o.p = &p
	return o
}

func NewPlayer(c net.Conn) *player {
	p := new(player)
	p.connection = c
	p.reader = bufio.NewReader(c)
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
	for _, o := range p.opponents { // this range should actually work like that
		sum = sum + o.scorevs
	}
	return sum
}

func (p *player) AddOpponent(newOpp player) opponent {
	o := NewOpponent(newOpp)
	p.opponents = append(p.opponents, *o)
	return *o
}

func (p1 *player) Challenge(p2 player) {
	fmt.Println(p1.name + " vs " + p2.name)
	p1op := p1.AddOpponent(p2)
	p2op := p2.AddOpponent(*p1)

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
