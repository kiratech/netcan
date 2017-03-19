package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/netcan/network"
	"github.com/fntlnz/netcan/util"
)

func main() {
	host, err := network.CreateHostFromPid("1")
	if err != nil {
		logrus.Fatal(err)
	}
	util.PrintHost(host)
}
