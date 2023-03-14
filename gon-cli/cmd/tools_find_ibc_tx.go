package cmd

import "github.com/spf13/cobra"

func findIBCTransactionsInteractive(cmd *cobra.Command) {
	sourceChain := chooseChain("Choose the source chain")
	destinationChain := chooseChain("Choose the destination chain")
	intialTxHash := askForString("Enter the transaction hash of the initial transaction")

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, intialTxHash)
}
