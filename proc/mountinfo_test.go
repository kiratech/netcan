package proc

import (
	"reflect"
	"testing"
)

type ParseMountData struct {
	rawline     string
	expectedset Mountinfo
}

// TestParseMountString data set, please add more cases if you feel
func ParseMountDataset() []ParseMountData {
	return []ParseMountData{
		{
			rawline: "515 24 0:3 net:[4026533140] /run/docker/netns/f46c0b2da189 rw shared:188 - nsfs nsfs rw",
			expectedset: Mountinfo{
				MountId:        "515",
				ParentId:       "24",
				MajorMinor:     "0:3",
				Root:           "net:[4026533140]",
				MountPoint:     "/run/docker/netns/f46c0b2da189",
				MountOptions:   "rw",
				OptionalFields: "shared:188",
				FilesystemType: "nsfs",
				MountSource:    "nsfs",
				SuperOptions:   "rw",
			},
		},
		{
			rawline: "26 25 0:24 / /sys/fs/cgroup/systemd rw,nosuid,nodev,noexec,relatime shared:9 - cgroup cgroup rw,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd",
			expectedset: Mountinfo{
				MountId:        "26",
				ParentId:       "25",
				MajorMinor:     "0:24",
				Root:           "/",
				MountPoint:     "/sys/fs/cgroup/systemd",
				MountOptions:   "rw,nosuid,nodev,noexec,relatime",
				OptionalFields: "shared:9",
				FilesystemType: "cgroup",
				MountSource:    "cgroup",
				SuperOptions:   "rw,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd",
			},
		},
	}
}

func TestParseMountString(t *testing.T) {
	for _, e := range ParseMountDataset() {
		info := ParseMountInfoString(e.rawline)

		if reflect.DeepEqual(e.expectedset, *info) == false {
			t.Error("Expected set is different than the resulting set")
		}
	}
}
