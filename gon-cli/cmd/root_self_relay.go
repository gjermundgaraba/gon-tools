package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
)

type Filter struct {
	chains.ChainData
}

func selfRelayInteractive(cmd *cobra.Command, args []string) {
	var (
		sourceChain      chains.Chain
		destinationChain chains.Chain
		txHash           string
	)

	if len(args) > 3 {
		sourceChainID := args[1]
		destinationChainID := args[2]
		txHash = args[3]

		if destinationChainID == sourceChainID {
			panic(fmt.Errorf("source and destination are the same chain"))
		}

		foundSource := false
		for _, chain := range chains.Chains {
			if string(chain.ChainID()) == sourceChainID {
				sourceChain = chain
				foundSource = true
				break
			}
		}

		if !foundSource {
			panic(fmt.Errorf("source chain %s not found", sourceChain))
		}

		foundDestination := false
		for _, chain := range chains.Chains {
			if string(chain.ChainID()) == destinationChainID {
				sourceChain = chain
				foundDestination = true
				break
			}
		}
		if !foundDestination {
			panic(fmt.Errorf("destination chain %s not found", destinationChainID))
		}
	} else {
		fmt.Println("This command requires the go relayer to have been set up according to the documentation see self-relay.md")
		youSure := askForConfirmation("This is currently an experimental feature, are you sure you want to continue?")
		if !youSure {
			fmt.Println("Alight! See you later :*")
			return
		}

		sourceChain = chooseChain("Source chain of transactions that needs relaying")
		destinationChain = chooseChain("Destination chain of transactions that needs relaying", sourceChain)

		txHash = askForString("Transaction hash to relay", survey.WithValidator(survey.Required))
	}

	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txHash, true, verbose)

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}
