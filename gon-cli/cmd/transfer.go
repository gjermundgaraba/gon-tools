package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func transferNFT(cmd *cobra.Command) error {
	sourceChain := chooseChain("Select source chain")
	setAddressPrefixes(sourceChain.Bech32Prefix())

	clientCtx := getClientContext(cmd, sourceChain)
	fromAccAddress := clientCtx.GetFromAddress()
	if fromAccAddress == nil {
		log.Fatal("No --from wallet/address specified")
	}
	fromAddress := sourceChain.ConvertAccAddressToChainsPrefix(fromAccAddress)

	destinationChain := chooseChain("Select destination chain", sourceChain)
	_ = destinationChain

	classes := sourceChain.ListNFTClassesThatHasNFTs(cmd.Context(), clientCtx, chains.ListNFTsQuery{
		Owner: fromAddress,
	})
	if len(classes) == 0 {
		log.Fatal("No NFT classes found")
	}

	selectedClass := chooseOne("Select class", classes)

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

	msg := sourceChain.CreateTransferNFTMsg(chosenChannel, selectedNFT, fromAddress, destinationAddress, timeoutHeight, timeoutTimestamp)
	if err := tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg); err != nil {
		return err
	}

	fmt.Println("Initial IBC transfer sent. It might take a moment before it is visible on the destination chain.")
	fmt.Println("The destination ibc trace will be:")
	fmt.Printf("%s/%s/%s\n", chosenConnection.ChannelB.Port, chosenConnection.ChannelB.Channel, selectedNFT.ClassID)

	return nil
}
