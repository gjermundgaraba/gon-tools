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
	"regexp"
)

var (
	// Taken from IRIS: https://github.com/irisnet/irismod/blob/main/modules/nft/types/validation.go
	// DenomID or TokenID can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number or a separator ('/', ':', '.', '_' or '-').
	idString = `[a-z][a-zA-Z0-9/]{2,127}`
	regexpID = regexp.MustCompile(fmt.Sprintf(`^%s$`, idString)).MatchString
)

func idValidator(val interface{}) error {
	// since we are validating an Input, the assertion will always succeed
	if str, ok := val.(string); !ok || !regexpID(str) {
		return fmt.Errorf("ClassID can only accept characters that match the regular expression: %s", idString)
	}
	return nil
}

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
	var chainOptions []chains.Chain
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

		chainOptions = append(chainOptions, chain)
	}

	return selectOne(questionPhrasing, chainOptions)
}

func chooseWallet(cmd *cobra.Command) string {
	if from, _ := cmd.Flags().GetString(flags.FlagFrom); from != "" {
		return from
	}

	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		panic(err)
	}

	records, err := clientCtx.Keyring.List()
	if err != nil {
		panic(err)
	}

	var walletNames []OptionString
	for _, o := range records {
		walletNames = append(walletNames, OptionString(o.Name))
	}

	return string(selectOne("Choose a wallet", walletNames))
}

type Option interface {
	Label() string
}

type OptionString string

func (o OptionString) Label() string {
	return string(o)
}

func selectOne[T Option](questionPhrasing string, options []T) T {
	var selectedIndex int
	var surveyOptions []string
	for _, o := range options {
		surveyOptions = append(surveyOptions, o.Label())
	}
	if err := survey.AskOne(&survey.Select{
		Message: questionPhrasing,
		Options: surveyOptions,
	}, &selectedIndex, survey.WithValidator(survey.Required)); err != nil {
		log.Fatalf("Error selecting: %v", err)
	}

	return options[selectedIndex]
}
