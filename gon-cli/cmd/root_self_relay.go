package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/gjermundgaraba/gon/gorelayer"
	"github.com/spf13/cobra"
)

type Filter struct {
	chains.ChainData
}

func selfRelayInteractive(cmd *cobra.Command, args []string, appHomeDir string) {
	var (
		sourceChain      chains.Chain
		destinationChain chains.Chain
		txHash           string
	)

	key := chooseOrCreateKey(cmd, chains.IRISChain)
	ethKey := getEthermintKeyName(key)
	verbose, _ := cmd.Flags().GetBool(flagVerbose)
	kr := getKeyring(cmd)
	rly := gorelayer.InitRly(appHomeDir, key, ethKey, kr, verbose)

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
		youSure := askForConfirmation("This is currently an experimental feature, are you sure you want to continue?", true)
		if !youSure {
			fmt.Println("Alight! See you later :*")
			return
		}

		sourceChain = chooseChain("Source chain of transactions that needs relaying")
		destinationChain = chooseChain("Destination chain of transactions that needs relaying", sourceChain)

		txHash = askForString("Transaction hash to relay", survey.WithValidator(survey.Required))
	}

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txHash, rly, true)

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}
