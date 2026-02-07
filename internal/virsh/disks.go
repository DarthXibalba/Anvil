package virsh

import (
	"fmt"
	"strings"
)

// Disk represents a VM block device.
type Disk struct {
	Target string // e.g. vda
	Source string // e.g. /var/lib/libvirt/images/vm1.qcow2
}

// ListDisks returns block devices for a domain using `virsh domblklist`.
func ListDisks(vm string) ([]Disk, error) {
	out, err := run("domblklist", vm, "--details")
	if err != nil {
		return nil, err
	}

	var disks []Disk
	lines := strings.Split(out, "\n")

	for _, l := range lines {
		fields := strings.Fields(l)
		if len(fields) < 4 {
			continue
		}

		// Expected format:
		// Type Device Target Source
		// disk file   vda    /path/to/image.qcow2
		if fields[0] == "disk" && fields[1] == "file" {
			disks = append(disks, Disk{
				Target: fields[2],
				Source: fields[3],
			})
		}
	}

	if len(disks) == 0 {
		return nil, fmt.Errorf("no file-backed disks found for VM %s", vm)
	}

	return disks, nil
}
