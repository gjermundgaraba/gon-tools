package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
)

func createNFTClassInteractive(cmd *cobra.Command) error {
	// TODO: Add support for juno, uptick
	chain := chooseChain("Select chain to create NFT on", chains.StargazeChain, chains.JunoChain, chains.UptickChain)
	setAddressPrefixes(chain.Bech32Prefix())

	key := chooseOrCreateKey(cmd, chain)
	if err := cmd.Flags().Set(flags.FlagFrom, key); err != nil {
		panic(err)
	}

	fromAddress := getAddressForChain(cmd, chain, key)

	var classID string
	if err := survey.AskOne(&survey.Input{Message: "Class ID"}, &classID, survey.WithValidator(idValidator)); err != nil {
		return err
	}

	var symbol string
	if err := survey.AskOne(&survey.Input{Message: "Symbol"}, &symbol, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var name string
	if err := survey.AskOne(&survey.Input{Message: "Name"}, &name, survey.WithValidator(survey.Required)); err != nil {
		return err
	}

	var description string
	if err := survey.AskOne(&survey.Input{Message: "Description"}, &description); err != nil {
		return err
	}

	var data string
	if err := survey.AskOne(&survey.Editor{
		Message:  "Data field (JSON)",
		FileName: "*.json",
	}, &data); err != nil {
		return err
	}

	var uri string
	if err := survey.AskOne(&survey.Input{Message: "URI"}, &uri); err != nil {
		return err
	}

	var uriHash string
	if uri != "" {
		if err := survey.AskOne(&survey.Input{Message: "URI Hash"}, &uriHash); err != nil {
			return err
		}
	}

	var mintRestricted bool
	var updateRestricted bool
	if chain.ChainID() == chains.IRISChain.ChainID() {
		if err := survey.AskOne(&survey.Confirm{Message: "Restrict minting NFTs to creator (you)?", Default: true}, &mintRestricted); err != nil {
			return err
		}

		if err := survey.AskOne(&survey.Confirm{Message: "Don't allow editing NFTs?", Default: true}, &updateRestricted); err != nil {
			return err
		}
	}

	msg := chain.CreateIssueCreditClassMsg(classID, name, "", fromAddress, symbol, mintRestricted, updateRestricted, description, uri, uriHash, data)

	clientCtx := getClientTxContext(cmd, chain)
	return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
}
