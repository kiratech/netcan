package network

import (
	"fmt"
	"net"

	"github.com/fntlnz/netcan/proc"
	"github.com/vishvananda/netlink"
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

func CreateHostFromPid(pid string) (*Host, error) {
	// Initialize everything
	netns := fmt.Sprintf("/proc/%s/ns/net", pid)
	mountinfo := fmt.Sprintf("/proc/%s/mountinfo", pid)

	interfaces := []*Interface{}
	hosts := []*Host{}
	links := []netlink.Link{}

	// Determine and append root infos
	rootIfaces, rootLinks, err := ExtractNetnsIfacesAndLinks(netns)

	if err != nil {
		return nil, fmt.Errorf("Error extracting network namespace and links from netns: %s => %s", netns, err)
	}
	rootHost := createHostFromRawIfaces(netns, rootIfaces)
	hosts = append(hosts, rootHost)
	links = append(links, rootLinks...)
	interfaces = append(interfaces, rootHost.Interfaces...)

	// Determine the sandboxes from mountinfo
	sandboxes, err := proc.GetMountInfo(mountinfo)
	if err != nil {
		return nil, err
	}

	// Extract info from all the sandboxes
	for _, c := range sandboxes {
		if c.FilesystemType != "nsfs" {
			continue
		}
		ifaces, curlinks, err := ExtractNetnsIfacesAndLinks(c.MountPoint)
		if err != nil {
			continue
		}
		curHost := createHostFromRawIfaces(c.MountPoint, ifaces)
		hosts = append(hosts, curHost)
		links = append(links, curlinks...)
		interfaces = append(interfaces, curHost.Interfaces...)
	}

	// Create associations
	flatIfaces := flattenInterfaces(interfaces)
	for _, l := range links {
		masterIdx := l.Attrs().MasterIndex
		parentIdx := l.Attrs().ParentIndex
		Idx := l.Attrs().Index

		for _, h := range hosts {
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

	return rootHost, nil
}
