package cmd

import (
	"fmt"
	"github.com/gjermundgaraba/gon/gorelayer"
	"log"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
)

func transferNFTInteractive(cmd *cobra.Command, appHomeDir string) error {
	sourceChain := chooseChain("Select source chain")
	setAddressPrefixes(sourceChain.Bech32Prefix())

	key := chooseOrCreateKey(cmd, sourceChain)
	if err := cmd.Flags().Set(flags.FlagFrom, key); err != nil {
		panic(err)
	}

	fromAddress := getAddressForChain(cmd, sourceChain, key)

	destinationChain := chooseChain("Select destination chain", sourceChain)
	_ = destinationChain

	clientCtx := getClientTxContext(cmd, sourceChain)
	selectedClass := getUsersNfts(cmd.Context(), clientCtx, sourceChain, fromAddress)
	if len(selectedClass.NFTs) == 0 {
		fmt.Println("No NFT classes found")
		return nil
	}

	// select nft
	// TODO: Use multiselect to be able to send more than one at a time
	selectedNFT := chooseOne("Select NFT", selectedClass.NFTs)

	var destinationAddress string
	if err := survey.AskOne(&survey.Input{Message: "What is the destination address? (Leave empty to send to same address as owner on destination)"}, &destinationAddress); err != nil {
		log.Fatalf("Error getting destination address: %v", err)
	}
	if destinationAddress == "" {
		destinationAddress = getAddressForChain(cmd, destinationChain, key)
		fmt.Println("Destination address:", destinationAddress)
	}

	chooseChannelQuestion := "Select channel to use"
	if selectedClass.LastIBCChannel.Port != "" {
		chooseChannelQuestion += fmt.Sprintf(" (last one was %s)", selectedClass.LastIBCChannel.Label())
	}
	chosenConnection := chooseConnection(sourceChain, destinationChain, chooseChannelQuestion)
	chosenChannel := chosenConnection.ChannelA

	tryToForceTimeout, _ := cmd.Flags().GetBool(flagTryToForceTimeout)
	targetChainHeight, targetChainTimestamp := getCurrentChainStatus(cmd.Context(), getQueryClientContext(cmd, destinationChain))
	timeoutHeight, timeoutTimestamp := sourceChain.GetIBCTimeouts(clientCtx, chosenChannel.Port, chosenChannel.Channel, targetChainHeight, targetChainTimestamp, tryToForceTimeout)

	msg := sourceChain.CreateTransferNFTMsg(chosenChannel, selectedClass, selectedNFT, fromAddress, destinationAddress, timeoutHeight, timeoutTimestamp)
	if tryToForceTimeout {
		clientCtx = clientCtx.WithSkipConfirmation(true)
	}
	txResponse, err := sendTX(clientCtx, cmd.Flags(), msg)
	if err != nil {
		panic(err)
	}

	fmt.Println("The destination ibc trace (full Class ID on destination chain) will be:")
	expectedDestinationClass, isRewind := calculateClassTrace(selectedClass.FullPathClassID, chosenConnection)
	if isRewind {
		fmt.Println("(This is a rewind transaction)")
	}
	fmt.Println(expectedDestinationClass)

	if len(strings.Split(expectedDestinationClass, "/")) > 2 && destinationChain.NFTImplementation() == chains.CosmosSDK {
		fmt.Println()
		fmt.Println("Class hash:")
		fmt.Println(calculateClassHash(expectedDestinationClass))
	}

	fmt.Println()
	selfRelay, _ := cmd.Flags().GetBool(flagSelfRelay)
	var rly *gorelayer.Rly
	if selfRelay {
		verbose, _ := cmd.Flags().GetBool(flagVerbose)
		kr := getKeyring(cmd)
		regularKeyName := getDefaultKeyName(key)
		ethKeyName := getEthermintKeyName(key)
		rly = gorelayer.InitRly(appHomeDir, regularKeyName, ethKeyName, kr, verbose)
	}

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txResponse.TxHash, rly, true)

	fmt.Println()
	fmt.Println("The destination ibc trace (full Class ID on destination chain):")
	if isRewind {
		fmt.Println("(This is a rewind transaction)")
	}
	fmt.Println(expectedDestinationClass)
	return nil
}
