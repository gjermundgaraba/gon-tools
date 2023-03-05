package cmd

import "github.com/spf13/cobra"

func CreateQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query NFTs",
		Long:  `Query for NFTs on any of the supported chains`,
	}

	cmd.AddCommand(CreateQueryClassesCmd())

	return cmd
}
