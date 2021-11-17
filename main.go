package main

import (
	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Type string      `json:"type"`
	User string      `json:"user"`
	Data interface{} `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		m := &Message{}
		err := conn.ReadJSON(m)
		if err != nil {
			log.Println(err)
			return
		}
		data := (*m).Data
		user := (*m).User

		(*m).Data = user + " : " + data.(string)

		err = conn.WriteJSON(m)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func main() {
	mux := pat.New()
	mux.Get("/chat", handler)
	n := negroni.Classic()
	n.UseHandler(mux)

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		return
	}
}
