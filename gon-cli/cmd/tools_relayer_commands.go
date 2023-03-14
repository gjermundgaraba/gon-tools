package cmd

import (
	"fmt"
	"strings"
)

var rlyPaths = []string{
	"gon-flixnet-1_channel-47-uni-6_channel-92",
	"gon-irishub-1_channel-1-gon-flixnet-1_channel-25",
	"gon-irishub-1_channel-19-uptick_7000-2_channel-4",
	"gon-irishub-1_channel-23-elgafar-1_channel-208",
	"uni-6_channel-94-elgafar-1_channel-213",
	"gon-flixnet-1_channel-46-uni-6_channel-91",
	"gon-irishub-1_channel-25-uni-6_channel-90",
	"uptick_7000-2_channel-13-uni-6_channel-88",
	"gon-flixnet-1_channel-44-elgafar-1_channel-209",
	"gon-flixnet-1_channel-45-elgafar-1_channel-210",
	"gon-irishub-1_channel-0-gon-flixnet-1_channel-24",
	"gon-irishub-1_channel-17-uptick_7000-2_channel-3",
	"gon-irishub-1_channel-24-uni-6_channel-89",
	"uni-6_channel-93-elgafar-1_channel-211",
	"uptick_7000-2_channel-5-gon-flixnet-1_channel-41",
	"uptick_7000-2_channel-7-uni-6_channel-86",
	"gon-irishub-1_channel-22-elgafar-1_channel-207",
	"uptick_7000-2_channel-6-elgafar-1_channel-203",
	"uptick_7000-2_channel-9-gon-flixnet-1_channel-42",
	"uptick_7000-2_channel-12-elgafar-1_channel-206",
}

func relayerCommandsInteractive() {
	sourceChain := chooseChain("Which chain is the source chain?")
	destinationChain := chooseChain("Which chain is the destination chain?", sourceChain)
	connection := chooseConnection(sourceChain, destinationChain, "Which connection do you want to relay?")

	rlySideA := fmt.Sprintf("%s_%s", sourceChain.ChainID(), connection.ChannelA.Channel)
	rlySideB := fmt.Sprintf("%s_%s", destinationChain.ChainID(), connection.ChannelB.Channel)
	fmt.Println(rlySideA)
	fmt.Println(rlySideB)
	var rlyPath string
	for _, rp := range rlyPaths {
		if strings.Contains(rp, rlySideA) && strings.Contains(rp, rlySideB) {
			rlyPath = rp
			break
		}
	}

	if rlyPath != "" {
		fmt.Println("To relay this connection with rly (go relayer), use the following command:")
		fmt.Printf("rly tx relay-packets %s %s\n", rlyPath, connection.ChannelA.Channel)
		fmt.Println()
	}

	fmt.Println("To relay this connection with hermes, use the following command:")
	fmt.Printf("hermes tx packet-recv --dst-chain %s --src-chain %s --src-port %s --src-channel %s\n", destinationChain.ChainID(), sourceChain.ChainID(), connection.ChannelA.Port, connection.ChannelA.Channel)

}
