package cmd

import (
	"fmt"

	"anvil/internal/virsh"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <vm>",
	Short: "Shutdown a VM gracefully",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := args[0]
		fmt.Printf("Stopping %s\n", vm)
		return virsh.Stop(vm)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
