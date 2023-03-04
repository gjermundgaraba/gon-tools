package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func setAddressPrefixes(prefix string) {
	accountPubKeyPrefix := prefix + "pub"
	validatorAddressPrefix := prefix + "valoper"
	validatorPubKeyPrefix := prefix + "valoperpub"
	consNodeAddressPrefix := prefix + "valcons"
	consNodePubKeyPrefix := prefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

func getClientContext(cmd *cobra.Command, chain chains.Chain) client.Context {
	if err := cmd.Flags().Set(flags.FlagNode, chain.RPC()); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set(flags.FlagGas, "auto"); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set(flags.FlagGasAdjustment, "1.5"); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set(flags.FlagGasPrices, fmt.Sprintf("0.25%s", chain.Denom())); err != nil {
		panic(err)
	}

	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		panic(err)
	}

	return clientCtx.WithChainID(string(chain.ChainID()))
}

func chooseChain(questionPhrasing string, filterOut ...chains.Chain) chains.Chain {
	var chainOptions []string
	for _, chain := range chains.Chains {
		toBeFilteredOut := false
		for _, filter := range filterOut {
			if chain.ChainID() == filter.ChainID() {
				toBeFilteredOut = true
				break
			}
		}
		if toBeFilteredOut {
			continue
		}

		chainOptions = append(chainOptions, chain.Name())
	}

	var selectedChainName string
	if err := survey.AskOne(&survey.Select{
		Message: questionPhrasing,
		Options: chainOptions,
	}, &selectedChainName, survey.WithValidator(survey.Required)); err != nil {
		log.Fatalf("Error selecting chain: %v", err)
	}

	for _, chain := range chains.Chains {
		if chain.Name() == selectedChainName {
			return chain
		}
	}

	panic("somehow didn't manage to choose chain")
}
