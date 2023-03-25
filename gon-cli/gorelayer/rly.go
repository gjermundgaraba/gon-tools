package gorelayer

import (
	"context"
	"fmt"
	rlycmd "github.com/cosmos/relayer/v2/cmd"
	"github.com/cosmos/relayer/v2/relayer"
	"github.com/cosmos/relayer/v2/relayer/provider"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var pathMap = map[string]map[string]string{
	"gon-flixnet-1": {
		"channel-24": "gon-irishub-1_channel-0-gon-flixnet-1_channel-24",
		"channel-25": "gon-irishub-1_channel-1-gon-flixnet-1_channel-25",
		"channel-41": "uptick_7000-2_channel-5-gon-flixnet-1_channel-41",
		"channel-42": "uptick_7000-2_channel-9-gon-flixnet-1_channel-42",
		"channel-44": "gon-flixnet-1_channel-44-elgafar-1_channel-209",
		"channel-45": "gon-flixnet-1_channel-45-elgafar-1_channel-210",
		"channel-46": "gon-flixnet-1_channel-46-uni-6_channel-91",
		"channel-47": "gon-flixnet-1_channel-47-uni-6_channel-92",
	},
	"gon-irishub-1": {
		"channel-0":  "gon-irishub-1_channel-0-gon-flixnet-1_channel-24",
		"channel-1":  "gon-irishub-1_channel-1-gon-flixnet-1_channel-25",
		"channel-17": "gon-irishub-1_channel-17-uptick_7000-2_channel-3",
		"channel-19": "gon-irishub-1_channel-19-uptick_7000-2_channel-4",
		"channel-22": "gon-irishub-1_channel-22-elgafar-1_channel-207",
		"channel-23": "gon-irishub-1_channel-23-elgafar-1_channel-208",
		"channel-24": "gon-irishub-1_channel-24-uni-6_channel-89",
		"channel-25": "gon-irishub-1_channel-25-uni-6_channel-90",
	},
	"uni-6": {
		"channel-86":  "uptick_7000-2_channel-7-uni-6_channel-86",
		"channel-88":  "uptick_7000-2_channel-13-uni-6_channel-88",
		"channel-89":  "gon-irishub-1_channel-24-uni-6_channel-89",
		"channel-90":  "gon-irishub-1_channel-25-uni-6_channel-90",
		"channel-91":  "gon-flixnet-1_channel-46-uni-6_channel-91",
		"channel-92":  "gon-flixnet-1_channel-47-uni-6_channel-92",
		"channel-93":  "uni-6_channel-93-elgafar-1_channel-211",
		"channel-94":  "uni-6_channel-94-elgafar-1_channel-213",
		"channel-120": "uni-6_channel-120-elgafar-1_channel-230",
		"channel-122": "uni-6_channel-122-elgafar-1_channel-234",
		"channel-133": "uni-6_elgafar-1_spam_test",
	},
	"uptick_7000-2": {
		"channel-3":  "gon-irishub-1_channel-17-uptick_7000-2_channel-3",
		"channel-4":  "gon-irishub-1_channel-19-uptick_7000-2_channel-4",
		"channel-5":  "uptick_7000-2_channel-5-gon-flixnet-1_channel-41",
		"channel-6":  "uptick_7000-2_channel-6-elgafar-1_channel-203",
		"channel-7":  "uptick_7000-2_channel-7-uni-6_channel-86",
		"channel-9":  "uptick_7000-2_channel-9-gon-flixnet-1_channel-42",
		"channel-12": "uptick_7000-2_channel-12-elgafar-1_channel-206",
		"channel-13": "uptick_7000-2_channel-13-uni-6_channel-88",
	},
	"elgafar-1": {
		"channel-203": "uptick_7000-2_channel-6-elgafar-1_channel-203",
		"channel-206": "uptick_7000-2_channel-12-elgafar-1_channel-206",
		"channel-207": "gon-irishub-1_channel-22-elgafar-1_channel-207",
		"channel-208": "gon-irishub-1_channel-23-elgafar-1_channel-208",
		"channel-209": "gon-flixnet-1_channel-44-elgafar-1_channel-209",
		"channel-210": "gon-flixnet-1_channel-45-elgafar-1_channel-210",
		"channel-211": "uni-6_channel-93-elgafar-1_channel-211",
		"channel-213": "uni-6_channel-94-elgafar-1_channel-213",
		"channel-230": "uni-6_channel-120-elgafar-1_channel-230",
		"channel-234": "uni-6_channel-122-elgafar-1_channel-234",
		"channel-241": "uni-6_elgafar-1_spam_test",
	},
}

type Rly struct {
	// Log is the root logger of the application.
	// Consumers are expected to store and use local copies of the logger
	// after modifying with the .With method.
	Log *zap.Logger

	Viper *viper.Viper

	HomePath string
	Debug    bool
	Config   *rlycmd.Config
}

func (rly *Rly) RelayPacket(ctx context.Context, connection chains.NFTConnection, packetSequence uint64) bool {
	src := rly.GetRelayerChain(string(connection.ChannelA.ChainID))
	dst := rly.GetRelayerChain(string(connection.ChannelB.ChainID))

	srch, dsth, err := relayer.QueryLatestHeights(ctx, src, dst)
	if err != nil {
		panic(err)
	}

	pathString := rly.GetPathString(connection)
	//var path *relayer.Path
	if _, err = rly.setPathsFromArgs(src, dst, pathString); err != nil {
		panic(err)
	}

	srcChannel, err := relayer.QueryChannel(ctx, src, connection.ChannelA.Channel)
	if err != nil {
		panic(err)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	var msgsSrc1, msgsDst1 []provider.RelayerMessage
	eg.Go(func() error {
		// Error ignored because it errors if there are no messages to relay, which might be fine! We deal with that later anyway
		err = relayer.AddMessagesForSequences(egCtx, []uint64{packetSequence}, src, dst, srch, dsth, &msgsSrc1, &msgsDst1,
			srcChannel.ChannelId, srcChannel.PortId, srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId, srcChannel.Ordering)
		rly.Log.Debug("Error from AddMessagesForSequences, msgsSrc1, msgsDst1", zap.String("src_chain_id", src.ChainID()), zap.String("dst_chain_id", dst.ChainID()), zap.Error(err))
		return nil
	})

	var msgsSrc2, msgsDst2 []provider.RelayerMessage
	eg.Go(func() error {
		// Error ignored because it errors if there are no messages to relay, which might be fine! We deal with that later anyway
		err = relayer.AddMessagesForSequences(egCtx, []uint64{packetSequence}, dst, src, dsth, srch, &msgsDst2, &msgsSrc2,
			srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId, srcChannel.ChannelId, srcChannel.PortId, srcChannel.Ordering)
		rly.Log.Debug("Error from AddMessagesForSequences, msgsSrc2, msgsDst2", zap.String("src_chain_id", dst.ChainID()), zap.String("dst_chain_id", src.ChainID()), zap.Error(err))
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
		rly.Log.Info(
			"No packets to relay",
			zap.String("src_chain_id", src.ChainID()),
			zap.String("src_port_id", srcChannel.PortId),
			zap.String("dst_chain_id", dst.ChainID()),
			zap.String("dst_port_id", srcChannel.Counterparty.PortId),
		)
		return false
	}

	if err := msgs.PrependMsgUpdateClient(ctx, src, dst, srch, dsth); err != nil {
		panic(err)
	}

	// send messages to their respective chains
	result := msgs.Send(ctx, rly.Log, relayer.AsRelayMsgSender(src), relayer.AsRelayMsgSender(dst), "self-relayed using the Game of NFTs CLI by @gjermundgaraba")
	if err := result.Error(); err != nil {
		if result.PartiallySent() {
			rly.Log.Info(
				"Partial success when relaying packets",
				zap.String("src_chain_id", src.ChainID()),
				zap.String("src_port_id", srcChannel.PortId),
				zap.String("dst_chain_id", dst.ChainID()),
				zap.String("dst_port_id", srcChannel.Counterparty.PortId),
				zap.Error(err),
			)
		}
		if err.Error() == "packet messages are redundant" {
			fmt.Println("Packet already relayed")
			return true
		}

		fmt.Println("Something wrong happened when relaying packets", err)
	}

	if acked := rly.RelayAcks(ctx, connection, packetSequence); acked {
		fmt.Println("ACKs relayed successfully")
	}

	return result.SuccessfulSrcBatches > 0 || result.SuccessfulDstBatches > 0
}

func (rly *Rly) RelayAcks(ctx context.Context, connection chains.NFTConnection, packetSequence uint64) bool {
	src := rly.GetRelayerChain(string(connection.ChannelA.ChainID))
	dst := rly.GetRelayerChain(string(connection.ChannelB.ChainID))
	srch, dsth, err := relayer.QueryLatestHeights(ctx, src, dst)
	if err != nil {
		panic(err)
	}

	srcChannel, err := relayer.QueryChannel(ctx, src, connection.ChannelA.Channel)
	if err != nil {
		panic(err)
	}

	// set the maximum relay transaction constraints
	msgs := &relayer.RelayMsgs{
		Src: []provider.RelayerMessage{},
		Dst: []provider.RelayerMessage{},
	}

	// add message for received packets on dst

	// dst wrote the ack. acknowledgementFromSequence will query the acknowledgement
	// from the counterparty chain (second chain provided in the arguments). The message
	// should be sent to src.

	// Error ignored because it errors if there are no messages to relay, which might be fine! We deal with that later anyway
	// TODO: Maybe print err in a verbose mode
	relayAckMsgs, _ := src.ChainProvider.AcknowledgementFromSequence(ctx, dst.ChainProvider, uint64(dsth), packetSequence, srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId, srcChannel.ChannelId, srcChannel.PortId)

	// Do not allow nil messages to the queued, or else we will panic in send()
	if relayAckMsgs != nil {
		msgs.Src = append(msgs.Src, relayAckMsgs)
	}

	// add messages for received packets on src

	// src wrote the ack. acknowledgementFromSequence will query the acknowledgement
	// from the counterparty chain (second chain provided in the arguments). The message
	// should be sent to dst.

	// Error ignored because it errors if there are no messages to relay, which might be fine! We deal with that later anyway
	// TODO: Maybe print err in a verbose mode
	relayAckMsgs, _ = dst.ChainProvider.AcknowledgementFromSequence(ctx, src.ChainProvider, uint64(srch), packetSequence, srcChannel.ChannelId, srcChannel.PortId, srcChannel.Counterparty.ChannelId, srcChannel.Counterparty.PortId)

	// Do not allow nil messages to the queued, or else we will panic in send()
	if relayAckMsgs != nil {
		msgs.Dst = append(msgs.Dst, relayAckMsgs)
	}

	if !msgs.Ready() {
		rly.Log.Info(
			"No acknowledgements to relay",
			zap.String("src_chain_id", src.ChainID()),
			zap.String("src_port_id", srcChannel.PortId),
			zap.String("dst_chain_id", dst.ChainID()),
			zap.String("dst_port_id", srcChannel.Counterparty.PortId),
		)
		return false
	}

	if err := msgs.PrependMsgUpdateClient(ctx, src, dst, srch, dsth); err != nil {
		panic(err)
	}

	// send messages to their respective chains
	result := msgs.Send(ctx, rly.Log, relayer.AsRelayMsgSender(src), relayer.AsRelayMsgSender(dst), "self-acked using the Game of NFTs CLI by @gjermundgaraba")
	if err := result.Error(); err != nil {
		if result.PartiallySent() {
			rly.Log.Info(
				"Partial success when relaying acknowledgements",
				zap.String("src_chain_id", src.ChainID()),
				zap.String("src_port_id", srcChannel.PortId),
				zap.String("dst_chain_id", dst.ChainID()),
				zap.String("dst_port_id", srcChannel.Counterparty.PortId),
				zap.Error(err),
			)
		}

		if err.Error() == "packet messages are redundant" {
			return true
		}

		panic(err)
	}

	return result.SuccessfulSrcBatches > 0 || result.SuccessfulDstBatches > 0
}

func (rly *Rly) GetRelayerChain(chainID string) *relayer.Chain {
	chain, err := rly.Config.Chains.Get(chainID)
	if err != nil {
		panic(err)
	}

	return chain
}

func (rly *Rly) GetPathString(connection chains.NFTConnection) string {
	path, ok := pathMap[string(connection.ChannelA.ChainID)][connection.ChannelA.Channel]
	if !ok {
		panic(fmt.Sprintf("Path not found for chainID %s and channel %s", connection.ChannelA.ChainID, connection.ChannelA.Channel))
	}

	return path
}

func (rly *Rly) setPathsFromArgs(src, dst *relayer.Chain, name string) (*relayer.Path, error) {
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
