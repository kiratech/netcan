package network

import (
	"fmt"

	"github.com/fntlnz/netcan/proc"
	"github.com/vishvananda/netlink"
)

type NetnsNetInfo struct {
	Interfaces []*Interface
	Hosts      []*Host
	Links      []netlink.Link
}

func AggregateNetnsNetworkInfo(netnsFd string, mountinfoFd string, rootfs string) (*NetnsNetInfo, error) {
	interfaces := []*Interface{}
	hosts := []*Host{}
	links := []netlink.Link{}

	// Determine and append root infos
	rootIfaces, rootLinks, err := extractNetnsIfacesAndLinks(netnsFd)

	if err != nil {
		return nil, fmt.Errorf("Error extracting network namespace and links from netns: %s => %s", netnsFd, err)
	}

	rootHost := createHostFromRawIfaces(netnsFd, rootIfaces)
	hosts = append(hosts, rootHost)
	links = append(links, rootLinks...)
	interfaces = append(interfaces, rootHost.Interfaces...)

	// Determine the sandboxes from mountinfo
	sandboxes, err := proc.GetMountInfo(mountinfoFd)
	if err != nil {
		return nil, err
	}

	// Extract info from all the sandboxes
	for _, c := range sandboxes {
		if c.FilesystemType != "nsfs" {
			continue
		}
		ifaces, curlinks, err := extractNetnsIfacesAndLinks(fmt.Sprintf("%s/%s", rootfs, c.MountPoint))
		if err != nil {
			continue
		}
		curHost := createHostFromRawIfaces(fmt.Sprintf("%s/%s", rootfs, c.MountPoint), ifaces)
		hosts = append(hosts, curHost)
		links = append(links, curlinks...)
		interfaces = append(interfaces, curHost.Interfaces...)
	}

	return &NetnsNetInfo{interfaces, hosts, links}, nil
}
