package network

import (
	"net"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

type VethPair struct {
	HostInterface   net.Interface
	GuestInterface  net.Interface
	BridgeInterface net.Interface
}

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

func GroupVethPairs(interfaces []net.Interface, links []netlink.Link) []VethPair {
	mapifaces := map[int]net.Interface{}

	for _, i := range interfaces {
		mapifaces[i.Index] = i
	}

	pairs := []VethPair{}
	for _, link := range links {
		masterIdx := link.Attrs().MasterIndex
		parentIdx := link.Attrs().ParentIndex
		Idx := link.Attrs().Index

		if _, ok := mapifaces[Idx]; !ok {
			continue
		}

		if _, ok := mapifaces[parentIdx]; !ok {
			continue
		}

		hostInterface := mapifaces[Idx]
		parentInterface := mapifaces[parentIdx]
		bridgeInterface := mapifaces[masterIdx]

		pairs = append(pairs, VethPair{
			HostInterface:   hostInterface,
			GuestInterface:  parentInterface,
			BridgeInterface: bridgeInterface,
		})
	}
	return pairs
}
