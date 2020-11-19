package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var tmpl *template.Template

type games struct {
	games    map[string][]string
	gamelens map[string]int
	sync.RWMutex
}

var currGames games

func main() {
	logFile, _ := os.Create("log")
	log.New(logFile, "", 1)

	currGames.games = make(map[string][]string)
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", createGame())
	http.HandleFunc("/add", addItem())
	http.HandleFunc("/play", play())
	http.HandleFunc("/nameRequest", nameRequest())
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func createGame() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var gameId string
		switch req.Method {
		case http.MethodGet:
			tmpl.ExecuteTemplate(w, "newGame", gameId)
			log.Printf("%s hoče igrati\n", req.RemoteAddr)
		case http.MethodPost:
			log.Printf("%s ustvarja igro\n", req.RemoteAddr)
			currGames.Lock()
			var ok bool
			for gameId = strconv.Itoa(rand.Int()); ok; _, ok = currGames.games[gameId] {
				gameId = strconv.Itoa(rand.Int())
				fmt.Println(gameId, ok)
			}
			currGames.games[gameId] = make([]string, 0, 10)
			currGames.Unlock()
			log.Printf("ustvaril igro %s\n", gameId)
			tmpl.ExecuteTemplate(w, "gameCreated", "/add?gameid="+gameId)
		}
	}
}

func addItem() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var gameId string
		switch req.Method {
		case http.MethodGet:
			gameId = req.FormValue("gameid")
			log.Printf("gameId %s\n", gameId)
			tmpl.ExecuteTemplate(w, "addItem", gameId)

		case http.MethodPost:
			gameId = req.FormValue("gameid")
			item := req.FormValue("item")
			log.Printf("%s v igri %s je dodal %s\n", req.Host, gameId, item)
			currGames.Lock()
			currGames.games[gameId] = append(currGames.games[gameId], item)
			a := currGames.games[gameId]
			rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
			currGames.Unlock()
			log.Println(currGames.games[gameId])
			http.Redirect(w, req, "/play?gameid="+gameId, 303)
		}
	}
}

func play() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var gameId string
		switch req.Method {
		case http.MethodGet:
			gameId = req.FormValue("gameid")
			log.Printf("%s v igri %s igra\n", req.Host, gameId)
			tmpl.ExecuteTemplate(w, "play", gameId)
		}
	}
}

func nameRequest() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		gameId := req.FormValue("gameid")
		fmt.Fprint(w, currGames.pop(gameId))
	}
}

func (g *games) pop(gameId string) string {
	if len(g.games[gameId]) <= 0 {
		return "zmanjkalo imen"
	}
	str := g.games[gameId][0]
	log.Print(g.games[gameId])
	g.games[gameId] = g.games[gameId][1:]
	return str
}
