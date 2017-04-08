package network

import "testing"

func TestCreateHostFromNetnsNetworkInfoWithoutHosts(t *testing.T) {
	ninfo := &NetnsNetInfo{}
	_, err := CreateHostFromNetnsNetworkInfo(ninfo)

	if err == nil {
		t.Error("Creating an hosts without any starting host should not happen")
	}
}
