package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func selfRelayInteractive(cmd *cobra.Command) {
	fmt.Println("This command requires the go relayer to have been set up according to the documentation see self-relay.md")
	youSure := askForConfirmation("This is currently an experimental feature, are you sure you want to continue?")
	if !youSure {
		fmt.Println("Alight! See you later :*")
		return
	}

	sourceChain := chooseChain("Source chain of transactions that needs relaying")
	destinationChain := chooseChain("Destination chain of transactions that needs relaying", sourceChain)

	txHash := askForString("Transaction hash to relay", survey.WithValidator(survey.Required))

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txHash, true)

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}
