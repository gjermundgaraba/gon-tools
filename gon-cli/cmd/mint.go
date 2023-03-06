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

	var classID string
	if err := survey.AskOne(&survey.Input{Message: "Class ID"}, &classID, survey.WithValidator(idValidator)); err != nil {
		return err
	}

	var nftID string
	if err := survey.AskOne(&survey.Input{Message: "NFT ID"}, &nftID, survey.WithValidator(idValidator)); err != nil {
		return err
	}

	var name string
	if err := survey.AskOne(&survey.Input{Message: "Name (Optional)"}, &name); err != nil {
		return err
	}

	var data string
	if err := survey.AskOne(&survey.Editor{
		Message:  "Data JSON (Optional)",
		FileName: "*.json",
	}, &data); err != nil {
		return err
	}

	var uri string
	if err := survey.AskOne(&survey.Input{Message: "URI (Optional)"}, &uri); err != nil {
		return err
	}

	var uriHash string
	if uri != "" {
		if err := survey.AskOne(&survey.Input{Message: "URI Hash (Optional)"}, &uriHash); err != nil {
			return err
		}
	}

	msg := chain.CreateMintNFTMsg(nftID, classID, name, uri, uriHash, data, fromAddress)
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
