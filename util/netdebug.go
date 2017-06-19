package util

import (
	"fmt"
	"strings"

	"github.com/kiratech/netcan/network"
)

func PrintIface(i *network.Interface, ind int) {
	indStr := strings.Repeat(" ", ind)
	fmt.Printf("%sName: %s\n", indStr, i.Name)
	fmt.Printf("%sIndex: %d\n", indStr, i.Index)

	if i.Pair != nil {
		fmt.Printf("%sPair {\n", indStr)
		if ind < 6 {
			printHost(i.Pair.Host, 2+ind)
		} else {
			fmt.Printf("%s%d - max depth reached\n", indStr, i.Pair.Index)
		}
		fmt.Printf("%s}\n", indStr)
	}
}

func printHost(h *network.Host, ind int) {
	indStr := strings.Repeat(" ", ind)
	fmt.Printf("%sNs: %s\n", indStr, h.Namespace.Fd)
	for _, i := range h.Interfaces {
		fmt.Println(fmt.Sprintf("%sInterface {", indStr))
		PrintIface(i, 2+ind)

		fmt.Println(fmt.Sprintf("%s  Bridges {", indStr))
		for _, b := range i.Bridges {
			PrintIface(b, 4+ind)
		}
		fmt.Println(fmt.Sprintf("%s  }", indStr))
		fmt.Println(fmt.Sprintf("%s}", indStr))
	}
}

func PrintHost(h *network.Host) {
	printHost(h, 0)
}
