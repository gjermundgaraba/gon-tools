package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/spf13/cobra"
)

func queryTransaction(cmd *cobra.Command) {
	chain := chooseChain("Select chain to create NFT on")
	setAddressPrefixes(chain.Bech32Prefix())

	clientCtx := getClientTxContext(cmd, chain)

	txHash := askForString("Transaction hash", survey.WithValidator(survey.Required))
	output, err := authtx.QueryTx(clientCtx, txHash)
	if err != nil {
		panic(err)
	}

	if output.Empty() {
		panic("No transaction found")
	}

	if err := clientCtx.PrintProto(output); err != nil {
		panic(err)
	}
}
