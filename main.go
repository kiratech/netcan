package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/netcan/network"
)

func main() {
	host, err := network.CreateHostFromPid("1")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info(host)
}
