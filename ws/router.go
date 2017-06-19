package ws

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/kiratech/netcan/network"
	"github.com/sirupsen/logrus"
)

const separator = "<->"

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

func checkIfEdgeExists(edges []csedge, edgeID string) bool {
	for _, ed := range edges {
		if ed.Data.Id == edgeID {
			return true
		}
		splitted := strings.Split(edgeID, separator)

		if len(splitted) < 2 {
			continue
		}
		reverse := formatPair(splitted[1], splitted[0])

		if ed.Data.Id == reverse {
			return true
		}

	}
	return false
}

func formatInterface(i network.Interface) string {
	return fmt.Sprintf("%d:%s", i.Index, i.Name)
}

func formatPair(left, right string) string {
	return fmt.Sprintf("%s%s%s", left, separator, right)
}

func appendHost(nodes []csnode, edges []csedge, host *network.Host) ([]csnode, []csedge) {
	// The parent host node
	nodes = append(nodes, csnode{
		csdata{
			Id: host.Namespace.Fd,
		},
	})

	for _, i := range host.Interfaces {
		interfaceStr := formatInterface(*i)
		nodes = append(nodes, csnode{
			csdata{
				Id:     interfaceStr,
				Parent: i.Host.Namespace.Fd,
			},
		})

		if i.Pair != nil {
			nodes, edges = appendInterface(nodes, edges, i.Pair)
		}

		for _, b := range i.Bridges {
			nodes = append(nodes, csnode{
				csdata{
					Id:     interfaceStr,
					Parent: i.Host.Namespace.Fd,
				},
			})

			curID := formatPair(interfaceStr, formatInterface(*b))

			if checkIfEdgeExists(edges, curID) {
				continue
			}
			edges = append(edges, csedge{
				csdata{
					Id:     curID,
					Source: formatInterface(*i),
					Target: formatInterface(*b),
					Style:  map[string]string{"line-color": "red"},
					Weight: 1,
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
			Id:     formatInterface(*i),
			Style:  map[string]string{"background-color": "lightsalmon"},
			Parent: i.Host.Namespace.Fd,
		},
	})

	if i.Pair != nil {
		curID := formatPair(formatInterface(*i), formatInterface(*i.Pair))
		if !checkIfEdgeExists(edges, curID) {
			edges = append(edges, csedge{
				csdata{
					Id:     curID,
					Source: formatInterface(*i),
					Target: formatInterface(*i.Pair),
					Weight: 1,
				},
			})
		}
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
