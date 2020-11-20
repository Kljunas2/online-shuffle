package main

import (
	"math/rand"
	"sync"
)

var errEmpty error

type game struct {
	sync.RWMutex
	id      string
	list    []string
	queue   []string
	players []string
}

type games struct {
	games map[string]*game
	sync.RWMutex
}

func createGame(id string) *game {
	newGame := &game{}
	newGame.id = id
	newGame.list = make([]string, 0, 16)
	newGame.players = make([]string, 0)
	return newGame
}

func (g *game) addItem(i string) {
	g.Lock()
	g.list = append(g.list, i)
	g.queue = append(g.queue, i)
	q := g.queue
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	g.Unlock()
}

func (g *game) pop() (string, error) {
	g.Lock()
	if len(g.queue) <= 0 {
		return "", errEmpty
	}
	i := g.queue[0]
	g.queue = g.queue[1:]
	g.Unlock()
	return i, nil
}

func (g *game) reshuffle() {
	g.Lock()
	copy(g.queue, g.list)
	q := g.queue
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	g.Unlock()
}
