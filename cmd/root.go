package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anvil",
	Short: "Anvil is a VM lifecycle CLI for libvirt-based labs",
	Long:  "Anvil manages KVM/libvirt virtual machines for fast cloning, lifecycle control, and cleanup.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
