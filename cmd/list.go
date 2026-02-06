package cmd

import (
	"fmt"

	"anvil/internal/virsh"

	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List VMs",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := virsh.List(listAll)
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "Show all VMs")
	rootCmd.AddCommand(listCmd)
}
