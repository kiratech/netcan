package ws

import (
	"encoding/json"
	"net/http"
	"time"

	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/fntlnz/netcan/network"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

func appendHost(nodes []csnode, edges []csedge, host *network.Host) ([]csnode, []csedge) {
	// The parent host node
	nodes = append(nodes, csnode{
		csdata{
			Id: host.Namespace.Fd,
		},
	})

	for _, i := range host.Interfaces {
		nodes = append(nodes, csnode{
			csdata{
				Id:     i.Name,
				Parent: i.Host.Namespace.Fd,
			},
		})

		if i.Pair != nil {
			nodes, edges = appendInterface(nodes, edges, i.Pair)
		}

		for _, b := range i.Bridges {
			nodes = append(nodes, csnode{
				csdata{
					Id:     i.Name,
					Parent: i.Host.Namespace.Fd,
				},
			})
			edges = append(edges, csedge{
				csdata{
					Id:        fmt.Sprintf("%s-%s", i.Name, b.Name),
					Source:    fmt.Sprintf("%s", i.Name),
					Target:    fmt.Sprintf("%s", b.Name),
					FaveShape: "triangle",
					Weight:    1,
				},
			})
		}
	}

	return nodes, edges
}

func appendInterface(nodes []csnode, edges []csedge, i *network.Interface) ([]csnode, []csedge) {
	// The parent node for this host
	nodes = append(nodes, csnode{
		csdata{
			Id: i.Host.Namespace.Fd,
		},
	})

	nodes = append(nodes, csnode{
		csdata{
			Id:     i.Name,
			Parent: i.Host.Namespace.Fd,
		},
	})

	if i.Pair != nil {
		edges = append(edges, csedge{
			csdata{
				Id:     fmt.Sprintf("%s-%s", i.Name, i.Pair.Name),
				Source: fmt.Sprintf("%s", i.Name),
				Target: fmt.Sprintf("%s", i.Pair.Name),
				Weight: 1,
			},
		})
	}

	return nodes, edges

}

func wsHandler(hub *Hub, rootfs string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
		go func() {
			var prevHash hash.Hash
			for {
				// TODO(fntlnz): worth waiting?
				time.Sleep(2 * time.Second)
				host, err := network.CreateHostFromPid("1", rootfs)
				if err != nil {
					logrus.Fatal(err)
				}
				cy := cytoscape{}
				nodes := []csnode{}
				edges := []csedge{}

				nodes, edges = appendHost(nodes, edges, host)

				cy.Elements = cselements{
					Nodes: nodes,
					Edges: edges,
				}

				j, err := json.Marshal(cy)
				if err != nil {
					logrus.Error(err)
					continue
				}

				curHash := sha1.New()
				curHash.Write(j)

				if prevHash != nil && hex.EncodeToString(curHash.Sum(nil)) == hex.EncodeToString(prevHash.Sum(nil)) {
					continue
				}
				hub.broadcast <- j
				prevHash = curHash
			}
		}()
	}
}

func NewRouter(rootfs string) *mux.Router {
	hub := newHub()
	go hub.run()
	r := mux.NewRouter()
	r.HandleFunc("/ws", wsHandler(hub, rootfs))
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("ui/"))))

	return r
}
