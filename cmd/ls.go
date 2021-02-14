package cmd

import "github.com/spf13/cobra"

func NewLSCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "list info",
	}
}
