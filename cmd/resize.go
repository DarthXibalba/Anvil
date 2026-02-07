package cmd

import (
	"fmt"

	"anvil/internal/qemu"
	"anvil/internal/virsh"

	"github.com/spf13/cobra"
)

var (
	resizeDisk        string
	resizeSize        string
	resizeAllowShrink bool
)

var resizeCmd = &cobra.Command{
	Use:   "resize <vm>",
	Short: "Resize a VM disk (qcow2)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := args[0]

		if resizeSize == "" {
			return fmt.Errorf("--size is required (e.g. 200G or +20G)")
		}

		disks, err := virsh.ListDisks(vm)
		if err != nil {
			return err
		}

		var target *virsh.Disk
		for _, d := range disks {
			if d.Target == resizeDisk {
				target = &d
				break
			}
		}

		if target == nil {
			return fmt.Errorf(
				"disk %q not found on VM %s (available: %s)",
				resizeDisk,
				vm,
				formatDiskTargets(disks),
			)
		}

		fmt.Printf(
			"Resizing %s:%s → %s\n",
			vm,
			target.Target,
			resizeSize,
		)

		return qemu.Resize(target.Source, resizeSize, resizeAllowShrink)
	},
}

func formatDiskTargets(disks []virsh.Disk) string {
	var out []string
	for _, d := range disks {
		out = append(out, d.Target)
	}
	return "[" + join(out, ", ") + "]"
}

func join(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	s := items[0]
	for i := 1; i < len(items); i++ {
		s += sep + items[i]
	}
	return s
}

func init() {
	resizeCmd.Flags().StringVar(&resizeDisk, "disk", "vda", "Disk target to resize (e.g. vda)")
	resizeCmd.Flags().StringVar(&resizeSize, "size", "", "New size (e.g. 200G or +20G)")
	resizeCmd.Flags().BoolVar(
		&resizeAllowShrink,
		"allow-shrink",
		false,
		"Allow shrinking the disk (DANGEROUS)",
	)

	rootCmd.AddCommand(resizeCmd)
}
