package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/gjermundgaraba/gon/gorelayer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Filter struct {
	chains.ChainData
}

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

func selfRelay(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	rly := gorelayer.InitRly(logger)

	sourceChainID := args[1]
	destinationChainID := args[2]
	txHash := args[3]

	var sourceChain chains.Chain
	var destinationChain chains.Chain

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

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txHash, true)

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}
