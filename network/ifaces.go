package network

import (
	"net"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func ExtractNetnsIfacesAndLinks(netnsfd string) ([]net.Interface, []netlink.Link, error) {
	ns, err := netns.GetFromPath(netnsfd)

	if err != nil {
		return nil, nil, err
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origns, _ := netns.Get()
	defer origns.Close()
	netns.Set(ns)
	defer ns.Close()

	ifaces, err := net.Interfaces()
	links, err := netlink.LinkList()

	netns.Set(origns)

	return ifaces, links, err
}
