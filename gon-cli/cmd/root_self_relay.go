package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	chantypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	"github.com/cosmos/relayer/v2/relayer"
	"github.com/cosmos/relayer/v2/relayer/provider"
	"github.com/gjermundgaraba/gon/gorelayer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

	src := rly.GetRelayerChain(sourceChain)
	dst := rly.GetRelayerChain(destinationChain)
	srch, dsth, err := relayer.QueryLatestHeights(ctx, src, dst)
	if err != nil {
		panic(err)
	}

	pathString := rly.GetPathString(connection)
	//var path *relayer.Path
	if _, err = setPathsFromArgs(rly, src, dst, pathString); err != nil {
		panic(err)
	}

	srcChannel, err := relayer.QueryChannel(cmd.Context(), src, connection.ChannelA.Channel)
	if err != nil {
		panic(err)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	var msgsSrc1, msgsDst1 []provider.RelayerMessage
	eg.Go(func() error {
		if err := relayer.AddMessagesForSequences(egCtx, []uint64{packetSequence}, src, dst, srch, dsth, &msgsSrc1, &msgsDst1,
			srcChannel.ChannelId, srcChannel.PortId, srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId, srcChannel.Ordering); err != nil {
			logger.Info(err.Error())
		}

		return nil
	})

	var msgsSrc2, msgsDst2 []provider.RelayerMessage
	eg.Go(func() error {
		if err := relayer.AddMessagesForSequences(egCtx, []uint64{packetSequence}, dst, src, dsth, srch, &msgsDst2, &msgsSrc2,
			srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId, srcChannel.ChannelId, srcChannel.PortId, srcChannel.Ordering); err != nil {
			logger.Info(err.Error())
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		panic(err)
	}

	msgs := &relayer.RelayMsgs{
		Src: append(msgsSrc1, msgsSrc2...),
		Dst: append(msgsDst1, msgsDst2...),
	}

	if !msgs.Ready() {
		logger.Info(
			"No packets to relay",
			zap.String("src_chain_id", src.ChainID()),
			zap.String("src_port_id", srcChannel.PortId),
			zap.String("dst_chain_id", dst.ChainID()),
			zap.String("dst_port_id", srcChannel.Counterparty.PortId),
		)
		return
	}

	if err := msgs.PrependMsgUpdateClient(ctx, src, dst, srch, dsth); err != nil {
		panic(err)
	}

	// send messages to their respective chains
	result := msgs.Send(ctx, logger, relayer.AsRelayMsgSender(src), relayer.AsRelayMsgSender(dst), "relayed using the Game of NFTs CLI by @gjermundgaraba")
	if err := result.Error(); err != nil {
		if result.PartiallySent() {
			logger.Info(
				"Partial success when relaying packets",
				zap.String("src_chain_id", src.ChainID()),
				zap.String("src_port_id", srcChannel.PortId),
				zap.String("dst_chain_id", dst.ChainID()),
				zap.String("dst_port_id", srcChannel.Counterparty.PortId),
				zap.Error(err),
			)
		}
		panic(err)
	}

	if result.SuccessfulSrcBatches > 0 {
		logPacketsRelayed(logger, src, dst, result.SuccessfulSrcBatches, srcChannel)
	}
	if result.SuccessfulDstBatches > 0 {
		logPacketsRelayed(logger, dst, src, result.SuccessfulDstBatches, srcChannel)
	}

	fmt.Println()
	fmt.Println("Relay seemingly successful!")
}

func logPacketsRelayed(logger *zap.Logger, src, dst *relayer.Chain, num int, srcChannel *chantypes.IdentifiedChannel) {
	logger.Info(
		"Relayed packets",
		zap.Int("count", num),
		zap.String("from_chain_id", dst.ChainID()),
		zap.String("from_port_id", srcChannel.Counterparty.PortId),
		zap.String("to_chain_id", src.ChainID()),
		zap.String("to_port_id", srcChannel.PortId),
	)
}

func setPathsFromArgs(rly *gorelayer.Rly, src, dst *relayer.Chain, name string) (*relayer.Path, error) {
	// find any configured paths between the chains
	paths, err := rly.Config.Paths.PathsFromChains(src.ChainID(), dst.ChainID())
	if err != nil {
		return nil, err
	}

	// Given the number of args and the number of paths, work on the appropriate
	// path.
	var path *relayer.Path
	switch {
	case name != "" && len(paths) > 1:
		if path, err = paths.Get(name); err != nil {
			return nil, err
		}

	case name != "" && len(paths) == 1:
		if path, err = paths.Get(name); err != nil {
			return nil, err
		}

	case name == "" && len(paths) > 1:
		return nil, fmt.Errorf("more than one path between %s and %s exists, pass in path name", src.ChainID(), dst.ChainID())

	case name == "" && len(paths) == 1:
		for _, v := range paths {
			path = v
		}
	}

	if err := src.SetPath(path.End(src.ChainID())); err != nil {
		return nil, err
	}

	if err := dst.SetPath(path.End(dst.ChainID())); err != nil {
		return nil, err
	}

	return path, nil
}
