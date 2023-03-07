package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func mintNFT(cmd *cobra.Command) error {
	// TODO: Add support for juno, uptick and omniflix
	chain := chooseChain("Select chain to mint NFT on", chains.StargazeChain, chains.JunoChain, chains.UptickChain, chains.OmniFlixChain)
	setAddressPrefixes(chain.Bech32Prefix())

	clientCtx := getClientContext(cmd, chain)
	fromAccAddress := clientCtx.GetFromAddress()
	if fromAccAddress == nil {
		log.Fatal("No --from wallet/address specified")
	}
	fromAddress := chain.ConvertAccAddressToChainsPrefix(fromAccAddress)

	classID := askForString("Class ID", survey.WithValidator(idValidator))

	nftID := askForString("NFT ID", survey.WithValidator(idValidator))

	name := askForString("Name (Optional)")

	var data string
	if err := survey.AskOne(&survey.Editor{
		Message:  "Data JSON (Optional)",
		FileName: "*.json",
	}, &data); err != nil {
		return err
	}

	uri := askForString("URI (Optional)")

	var uriHash string
	if uri != "" {
		uriHash = askForString("URI Hash (Optional)")
	}

	msg := chain.CreateMintNFTMsg(nftID, classID, name, uri, uriHash, data, fromAddress)
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
