package ws

import (
	"fmt"
	"net/http"
	"github.com/fntlnz/netcan/network"
	"github.com/Sirupsen/logrus"
)

func GraphHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	w.WriteHeader(http.StatusOK)
	host, err := network.CreateHostFromPid("1")
	if err != nil {
		logrus.Fatal(err)
	}

	w.Write([]byte(fmt.Sprintf("prova: %s", host.Namespace)))


}