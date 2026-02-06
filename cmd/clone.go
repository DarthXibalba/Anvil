package cmd

import (
	"fmt"

	"anvil/internal/virsh"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source-vm> <new-vm>",
	Short: "Clone a VM from an existing domain",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src := args[0]
		dst := args[1]

		fmt.Printf("Cloning %s -> %s\n", src, dst)
		return virsh.Clone(src, dst)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
