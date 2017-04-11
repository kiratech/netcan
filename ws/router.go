package ws

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/fntlnz/netcan/network"
	"fmt"
)

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	go client.writePump()
	client.readPump()
}

func NewRouter() *mux.Router {
	hub := newHub()
	go hub.run()
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	r.Handle("/", r)

	go func () {
		for {
			host, err := network.CreateHostFromPid("1")
			if err != nil {
				logrus.Fatal(err)
			}
			hub.broadcast <- []byte(fmt.Sprintf("%s", host.Namespace))
		}
	}()
	return r
}
