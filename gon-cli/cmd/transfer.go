package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

func transferNFT(cmd *cobra.Command) error {
	sourceChain := chooseChain("Select source chain", chains.UptickChain)
	setAddressPrefixes(sourceChain.Bech32Prefix())

	clientCtx := getClientContext(cmd, sourceChain)
	fromAccAddress := clientCtx.GetFromAddress()
	if fromAccAddress == nil {
		log.Fatal("No --from wallet/address specified")
	}
	fromAddress := sourceChain.ConvertAccAddressToChainsPrefix(fromAccAddress)

	destinationChain := chooseChain("Select destination chain", sourceChain, chains.UptickChain)
	_ = destinationChain

	selectedClass := getUsersNfts(cmd.Context(), clientCtx, sourceChain, fromAddress)
	if len(selectedClass.NFTs) == 0 {
		panic("No NFT classes found")
	}

	// select nft
	// TODO: Use multiselect to be able to send more than one at a time
	selectedNFT := chooseOne("Select NFT", selectedClass.NFTs)

	var destinationAddress string
	if err := survey.AskOne(&survey.Input{Message: "What is the destination address? (Leave empty to send to same address as owner on destination)"}, &destinationAddress); err != nil {
		log.Fatalf("Error getting destination address: %v", err)
	}
	if destinationAddress == "" {
		destinationAddress = destinationChain.ConvertAddressToChainsPrefix(fromAddress)
	}

	connections := sourceChain.GetConnectionsTo(destinationChain)
	var wrappedConnections []OptionWrapper[chains.NFTConnection]
	for _, connection := range connections {
		wrappedConnections = append(wrappedConnections, OptionWrapper[chains.NFTConnection]{
			WrappedValue: connection,
			LabelFunc: func(connection chains.NFTConnection) string {
				return connection.ChannelA.Label()
			},
		})
	}

	chooseChannelQuestion := "Select channel to use"
	if selectedClass.LastIBCChannel.Port != "" {
		chooseChannelQuestion += fmt.Sprintf(" (last one was %s)", selectedClass.LastIBCChannel.Label())
	}
	chosenConnection := chooseOne(chooseChannelQuestion, wrappedConnections).WrappedValue
	chosenChannel := chosenConnection.ChannelA

	timeoutHeight, timeoutTimestamp := sourceChain.GetIBCTimeouts(clientCtx, chosenChannel.Port, chosenChannel.Channel)

	msg := sourceChain.CreateTransferNFTMsg(chosenChannel, selectedClass, selectedNFT, fromAddress, destinationAddress, timeoutHeight, timeoutTimestamp)
	if err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg); err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Println("Initial IBC transfer sent. It might take a moment before it is visible on the destination chain.")
	fmt.Println()
	fmt.Println("The destination ibc trace (full Class ID on destination chain) will be:")
	var expectedDestinationClass string
	if strings.HasPrefix(selectedClass.FullPathClassID, fmt.Sprintf("%s/%s", chosenConnection.ChannelA.Port, chosenConnection.ChannelA.Channel)) {
		fmt.Println("(This is a rewind transaction)")
		expectedDestinationClass = strings.TrimPrefix(selectedClass.FullPathClassID, fmt.Sprintf("%s/%s/", chosenConnection.ChannelA.Port, chosenConnection.ChannelA.Channel))
	} else {
		expectedDestinationClass = fmt.Sprintf("%s/%s/%s", chosenConnection.ChannelB.Port, chosenConnection.ChannelB.Channel, selectedClass.FullPathClassID)
	}
	fmt.Println(expectedDestinationClass)

	return nil
}
