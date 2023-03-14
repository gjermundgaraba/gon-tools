package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gjermundgaraba/gon/gorelayer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"strconv"
)

func selfRelayInteractive(cmd *cobra.Command) {
	ctx := cmd.Context()
	fmt.Println("This command requires the go relayer to have been set up according to the documentation see self-relay.md")
	youSure := askForConfirmation("This is currently an experimental feature, are you sure you want to continue?")
	if !youSure {
		fmt.Println("Alight! See you later :*")
		return
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	rly := gorelayer.InitRly(logger)

	sourceChain := chooseChain("Source chain of transactions that needs relaying")
	destinationChain := chooseChain("Destination chain of transactions that needs relaying", sourceChain)

	txHash := askForString("Transaction hash to relay", survey.WithValidator(survey.Required))

	txResp := waitForTX(cmd, sourceChain, txHash, "Initial IBC packet", "Initial IBC packet")
	packetSequence, err := strconv.ParseUint(findPacketSequence(txResp), 10, 64)
	if err != nil {
		panic(err)
	}
	connection := findConnection(txResp)
	connection.ChannelA.ChainID = sourceChain.ChainID()
	connection.ChannelB.ChainID = destinationChain.ChainID()

	rly.RelayPacket(ctx, connection, packetSequence)

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}
