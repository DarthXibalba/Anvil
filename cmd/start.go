package cmd

import (
	"fmt"

	"anvil/internal/virsh"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <vm>",
	Short: "Start a VM",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vm := args[0]
		fmt.Printf("Starting %s\n", vm)
		return virsh.Start(vm)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
