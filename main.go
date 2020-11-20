package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var port string

func init() {
	flag.StringVar(&port, "port", ":8000", "port to listen")
	flag.Parse()
}

var tmpl *template.Template

var currGames games

func main() {
	currGames.games = make(map[string]*game)
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/", mainPage())
	http.HandleFunc("/add", addItem())
	http.HandleFunc("/play", play())
	http.HandleFunc("/nameRequest", nameRequest())

	staticServer := http.FileServer(http.Dir("."))
	http.Handle("/static/", staticServer)

	log.Fatal(http.ListenAndServe(port, nil))
}

func mainPage() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var gameId string
		switch req.Method {
		case http.MethodGet:
			tmpl.ExecuteTemplate(w, "newGame", gameId)
			log.Printf("%s hoče igrati\n", req.RemoteAddr)
		case http.MethodPost:
			fmt.Println(req.FormValue("gameid"))
			fmt.Println(req.FormValue("create"))
			fmt.Println(req.FormValue("play"))
			log.Printf("%s ustvarja igro\n", req.RemoteAddr)
			currGames.Lock()
			var ok bool
			for gameId = randhex(); ok; _, ok = currGames.games[gameId] {
				gameId = randhex()
				fmt.Println(gameId, ok)
			}
			currGames.games[gameId] = createGame(gameId)
			currGames.Unlock()
			log.Printf("ustvaril igro %s\n", gameId)
			tmpl.ExecuteTemplate(w, "gameCreated", "/add?gameid="+gameId)
		}
	}
}

func randhex() string {
	return fmt.Sprintf("%05x", rand.Intn(0x100000))
}

func addItem() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var gameId string
		switch req.Method {
		case http.MethodGet:
			gameId = req.FormValue("gameid")
			tmpl.ExecuteTemplate(w, "addItem", gameId)
		case http.MethodPost:
			gameId = req.FormValue("gameid")
			item := req.FormValue("item")
			log.Printf("%s v igri %s je dodal %s\n", req.Host, gameId, item)
			currGames.games[gameId].addItem(item)
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
	if len(g.games[gameId].list) <= 0 {
		return "zmanjkalo imen"
	}
	str := g.games[gameId].list[0]
	log.Print(g.games[gameId])
	g.games[gameId].list = g.games[gameId].list[1:]
	return str
}

func loginCookie(username string) http.Cookie {
	cookieValue := username + ":" /*+ codify.SHA(username+strconv.Itoa(rand.Intn(10000000)))*/
	expire := time.Now().AddDate(0, 0, 10)
	return http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire}
}
