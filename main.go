package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const default_netns_path = "/proc/1/ns/net"

type VethPair struct {
	HostInterface  net.Interface
	GuestInterface net.Interface
}

type NamespacePaths struct {
	Newnet string `json:"NEWNET"`
}

type OCIContainer struct {
	Id             string         `json:"id"`
	NamespacePaths NamespacePaths `json:"namespace_paths"`
}

func getContainer(dir os.FileInfo) (*OCIContainer, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("/run/runc/%s/state.json", dir.Name()))
	if err != nil {
		return nil, err
	}

	ocic := &OCIContainer{}
	err = json.Unmarshal(file, &ocic)

	if err != nil {
		return nil, err
	}

	return ocic, nil
}

func getContainers() ([]OCIContainer, error) {
	var containers []OCIContainer
	list, err := ioutil.ReadDir("/run/runc")

	if err != nil {
		return containers, err
	}

	for _, dir := range list {

		if dir.IsDir() == false {
			continue
		}

		ocic, err := getContainer(dir)

		if err != nil {
			logrus.Info("error getting container: ", err.Error())
			continue
		}

		containers = append(containers, *ocic)
	}

	return containers, nil
}

func extractIfaces(netnsfd string) ([]net.Interface, error) {
	ns, err := netns.GetFromPath(netnsfd)

	if err != nil {
		return nil, err
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Save the current network namespace
	origns, _ := netns.Get()
	defer origns.Close()
	netns.Set(ns)
	defer ns.Close()

	ifaces, err := net.Interfaces()

	netns.Set(origns)

	return ifaces, err
}

func printIfaces(ifaces []net.Interface) {
	for _, iface := range ifaces {
		logrus.Infof("Index: %d\tName: %s\tFlags:%v", iface.Index, iface.Name, iface.Flags)
	}
}

func printLink(link netlink.Link) {
	logrus.Info("encap: ", link.Attrs().EncapType, " Masteridx: ", link.Attrs().MasterIndex, " Parentidx: ", link.Attrs().ParentIndex, " Idx: ", link.Attrs().Index)
}

func groupPairs(interfaces []net.Interface, links []netlink.Link) []VethPair {
	mapifaces := map[int]net.Interface{}

	for _, i := range interfaces {
		mapifaces[i.Index] = i
	}

	pairs := []VethPair{}
	for _, link := range links {
		//masterIdx := link.Attrs().MasterIndex
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

		pairs = append(pairs, VethPair{
			HostInterface:  hostInterface,
			GuestInterface: parentInterface,
		})
	}
	return pairs
}

func main() {
	containers, err := getContainers()
	if err != nil {
		logrus.Fatal(err)
	}

	interfaces := []net.Interface{}

	logrus.Infof("HOST NS: %s", default_netns_path)
	hifaces, _ := extractIfaces(default_netns_path)
	interfaces = append(interfaces, hifaces...)

	for _, c := range containers {
		logrus.Infof("NS: %s", c.NamespacePaths.Newnet)
		ifaces, _ := extractIfaces(c.NamespacePaths.Newnet)
		interfaces = append(interfaces, ifaces...)
	}

	links, _ := netlink.LinkList()

	pairs := groupPairs(interfaces, links)

	logrus.Infof("Pairs: %v", pairs)
}
