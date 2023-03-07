package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func listConnections(cmd *cobra.Command) error {
	sourceChain := chooseChain("From chain:")
	destinationChain := chooseChain("To chain:")

	connections := sourceChain.GetConnectionsTo(destinationChain)
	for i, connection := range connections {
		fmt.Printf("Connection %d:\n", i+1)
		fmt.Printf("Source port: %s\n", connection.ChannelA.Port)
		fmt.Printf("Source channel: %s\n", connection.ChannelA.Channel)
		fmt.Printf("Destination port: %s\n", connection.ChannelB.Port)
		fmt.Printf("Destination channel: %s\n", connection.ChannelB.Channel)
		if i != len(connections)-1 {
			fmt.Println()
		}
	}

	return nil
}
