package util

import (
	"fmt"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/netcan/network"
	"github.com/vishvananda/netlink"
)

func PrintPairs(pairs []network.VethPair) {
	for _, pair := range pairs {
		fmt.Printf("%s=>%d <--%s---> %s=>%d\n", pair.GuestInterface.Name, pair.GuestInterface.Index, pair.BridgeInterface.Name, pair.HostInterface.Name, pair.HostInterface.Index)
	}
}

func PrintIfaces(ifaces []net.Interface) {
	for _, iface := range ifaces {
		logrus.Infof("Index: %d\tName: %s\tFlags:%v", iface.Index, iface.Name, iface.Flags)
	}
}

func PrintLink(link netlink.Link) {
	logrus.Info("encap: ", link.Attrs().EncapType, " Masteridx: ", link.Attrs().MasterIndex, " Parentidx: ", link.Attrs().ParentIndex, " Idx: ", link.Attrs().Index)
}
