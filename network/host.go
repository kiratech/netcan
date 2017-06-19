package network

import (
	"fmt"
	"net"
)

type Namespace struct {
	Fd string
}

type Interface struct {
	net.Interface
	Pair    *Interface
	Bridges []*Interface
	Host    *Host
}

type Host struct {
	Namespace  Namespace
	Interfaces []*Interface
}

func CreateHost(netns string, interfaces []*Interface) *Host {
	ns := Namespace{Fd: netns}
	return &Host{
		Namespace:  ns,
		Interfaces: interfaces,
	}
}

func createHostFromRawIfaces(netns string, ifaces []net.Interface) *Host {
	host := CreateHost(netns, nil)
	interfaces := []*Interface{}
	for _, iface := range ifaces {
		bridges := []*Interface{}
		interfaces = append(interfaces, &Interface{
			Interface: iface,
			Host:      host,
			Bridges:   bridges,
		})
	}

	host.Interfaces = interfaces
	return host
}

func flattenInterfaces(ifaces []*Interface) map[int]*Interface {
	interfaces := map[int]*Interface{}
	for _, i := range ifaces {
		interfaces[i.Index] = i
	}
	return interfaces
}

func CreateHostFromPid(pid string, rootfs string) (*Host, error) {
	netns := fmt.Sprintf("%s/proc/%s/ns/net", rootfs, pid)
	mountinfo := fmt.Sprintf("%s/proc/%s/mountinfo", rootfs, pid)
	return CreateHostFromPaths(netns, mountinfo, rootfs)
}

func CreateHostFromPaths(netns string, mountinfo string, rootfs string) (*Host, error) {
	netnsNetworkInfo, err := AggregateNetnsNetworkInfo(netns, mountinfo, rootfs)
	if err != nil {
		return nil, err
	}

	return CreateHostFromNetnsNetworkInfo(netnsNetworkInfo)
}

func CreateHostFromNetnsNetworkInfo(netnsNetworkInfo *NetnsNetInfo) (*Host, error) {

	if len(netnsNetworkInfo.Hosts) < 1 {
		return nil, fmt.Errorf("Unable to create an host given the provided paths")
	}

	// Create associations
	flatIfaces := flattenInterfaces(netnsNetworkInfo.Interfaces)

	for _, l := range netnsNetworkInfo.Links {
		masterIdx := l.Attrs().MasterIndex
		parentIdx := l.Attrs().ParentIndex
		Idx := l.Attrs().Index

		for _, h := range netnsNetworkInfo.Hosts {
			for _, i := range h.Interfaces {
				if i.Index != Idx {
					continue
				}
				if pairIf, ok := flatIfaces[parentIdx]; ok {
					i.Pair = pairIf
				}
				if brif, ok := flatIfaces[masterIdx]; ok {
					i.Bridges = append(i.Bridges, brif)
				}
			}
		}
	}

	return netnsNetworkInfo.Hosts[0], nil
}
