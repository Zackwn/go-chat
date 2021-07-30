package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveHome(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, req, "home.html")
}

func serveWs(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		s.newClient(ws)
	}
}

func main() {
	s := NewServer()
	go s.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs(s))
	log.Println("ready")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
