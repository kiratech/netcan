package main

import (
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/fntlnz/netcan/network"
	"github.com/fntlnz/netcan/proc"
	"github.com/fntlnz/netcan/util"
	"github.com/vishvananda/netlink"
)

const default_netns_path = "/proc/1/ns/net"

func main() {
	sandboxes, err := proc.GetMountInfo("/proc/1/mountinfo")

	if err != nil {
		logrus.Fatal(err)
	}

	interfaces := []net.Interface{}
	links := []netlink.Link{}

	logrus.Infof("HOST NS: %s", default_netns_path)
	hifaces, hlinks, _ := network.ExtractNetnsIfacesAndLinks(default_netns_path)
	interfaces = append(interfaces, hifaces...)
	links = append(links, hlinks...)

	for _, c := range sandboxes {
		if c.FilesystemType != "nsfs" {
			continue
		}
		logrus.Infof("NS: %s", c.MountPoint)
		ifaces, curlinks, _ := network.ExtractNetnsIfacesAndLinks(c.MountPoint)
		util.PrintIfaces(ifaces)
		interfaces = append(interfaces, ifaces...)
		links = append(links, curlinks...)
	}

	pairs := network.GroupVethPairs(interfaces, links)

	util.PrintPairs(pairs)
}
