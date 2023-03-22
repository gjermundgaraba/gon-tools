package cmd

import "github.com/spf13/cobra"

func findIBCTransactionsInteractive(cmd *cobra.Command) {
	sourceChain := chooseChain("Choose the source chain")
	destinationChain := chooseChain("Choose the destination chain", sourceChain)
	intialTxHash := askForString("Enter the transaction hash of the initial transaction")

	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, intialTxHash, false, verbose, true)
}
