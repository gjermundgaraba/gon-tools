package cmd

import (
	"fmt"
	"strings"

	"github.com/gjermundgaraba/gon/chains"
)

var (
	ReverseChainMap = map[chains.ChainID]string{
		chains.IRISChain.ChainID():     "i",
		chains.StargazeChain.ChainID(): "s",
		chains.JunoChain.ChainID():     "j",
		chains.OmniFlixChain.ChainID(): "o",
		chains.UptickChain.ChainID():   "u",
	}
)

func listConnections() {
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
}

func getMatrix() map[string]string {
	matrix := make(map[string]string)
	connections := chains.Connections
	var prevChainA chains.ChainID
	var prevChainB chains.ChainID
	for _, connection := range connections {
		l1 := ReverseChainMap[connection.ChannelA.ChainID]
		l2 := ReverseChainMap[connection.ChannelB.ChainID]
		if prevChainA == connection.ChannelA.ChainID && prevChainB == connection.ChannelB.ChainID {
			matrix[strings.ToUpper(l1)+l2] = connection.ChannelB.Port + "/" + connection.ChannelB.Channel
			matrix[strings.ToUpper(l1+l2)] = connection.ChannelB.Port + "/" + connection.ChannelB.Channel
			matrix[l2+l1] = connection.ChannelA.Port + "/" + connection.ChannelA.Channel
			matrix[strings.ToUpper(l2)+l1] = connection.ChannelA.Port + "/" + connection.ChannelA.Channel
		} else {
			matrix[l1+l2] = connection.ChannelB.Port + "/" + connection.ChannelB.Channel
			matrix[l1+strings.ToUpper(l2)] = connection.ChannelB.Port + "/" + connection.ChannelB.Channel
			matrix[l2+l1] = connection.ChannelA.Port + "/" + connection.ChannelA.Channel
			matrix[l2+strings.ToUpper(l1)] = connection.ChannelA.Port + "/" + connection.ChannelA.Channel

		}
		prevChainA = connection.ChannelA.ChainID
		prevChainB = connection.ChannelB.ChainID
	}
	return matrix
}
