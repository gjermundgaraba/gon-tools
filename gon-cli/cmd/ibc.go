package cmd

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	"strconv"
)

func findIBCTransactionsInteractive(cmd *cobra.Command) {
	sourceChain := chooseChain("Choose the source chain")
	destinationChain := chooseChain("Choose the destination chain")
	intialTxHash := askForString("Enter the transaction hash of the initial transaction")

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, intialTxHash)
}

func waitAndPrintIBCTrail(cmd *cobra.Command, sourceChain chains.Chain, destinationChain chains.Chain, txHash string) {
	txResp := waitForTX(cmd, sourceChain, txHash, "Initial IBC packet", "Initial IBC packet")
	packetSequence := findPacketSequence(txResp)
	connection := findConnection(txResp)
	timeoutHeight, timeoutTimestamp := findTimeouts(txResp)
	waitForIBCPacket(cmd, sourceChain, destinationChain, connection, packetSequence, timeoutHeight, timeoutTimestamp)
}

func findPacketSequence(txResp *sdk.TxResponse) string {
	for _, event := range txResp.Events {
		if event.Type == "send_packet" {
			for _, attr := range event.Attributes {
				if string(attr.Key) == "packet_sequence" {
					return string(attr.Value)
				}
			}
		}
	}
	panic("No packet sequence found")
}

func findConnection(txResp *sdk.TxResponse) chains.NFTConnection {
	var connection chains.NFTConnection

	for _, event := range txResp.Events {
		if event.Type == "send_packet" || event.Type == "recv_packet" || event.Type == "timeout_packet" {
			for _, attr := range event.Attributes {
				switch string(attr.Key) {
				case "packet_src_port":
					connection.ChannelA.Port = string(attr.Value)
				case "packet_src_channel":
					connection.ChannelA.Channel = string(attr.Value)
				case "packet_dst_port":
					connection.ChannelB.Port = string(attr.Value)
				case "packet_dst_channel":
					connection.ChannelB.Channel = string(attr.Value)
				}
			}
		}
	}

	if connection.ChannelA.Port == "" || connection.ChannelA.Channel == "" || connection.ChannelB.Port == "" || connection.ChannelB.Channel == "" {
		panic("No connection found")
	}

	return connection
}

func findTimeouts(txResp *sdk.TxResponse) (timeoutHeight uint64, timeoutTimestamp uint64) {
	foundEither := false
	for _, event := range txResp.Events {
		if event.Type == "send_packet" || event.Type == "recv_packet" || event.Type == "timeout_packet" {
			for _, attr := range event.Attributes {
				switch string(attr.Key) {
				case "packet_timeout_height":
					height, err := clienttypes.ParseHeight(string(attr.Value))
					if err != nil {
						panic(err)
					}
					timeoutHeight = height.GetRevisionHeight()
					foundEither = true
				case "packet_timeout_timestamp":
					var err error
					timeoutTimestamp, err = strconv.ParseUint(string(attr.Value), 10, 64)
					if err != nil {
						panic(err)
					}
					foundEither = true
				}
			}
		}
	}

	if !foundEither {
		panic("No timeouts found")
	}

	return timeoutHeight, timeoutTimestamp
}

func waitForIBCPacket(cmd *cobra.Command, sourceChain chains.Chain, destinationChain chains.Chain, connection chains.NFTConnection, packetSequence string, timeoutHeight uint64, timeoutTimestamp uint64) {
	_, timeout := waitForTXByEvents(
		cmd,
		destinationChain,
		[]string{
			fmt.Sprintf("recv_packet.packet_sequence='%s'", packetSequence),
			fmt.Sprintf("recv_packet.packet_src_port='%s'", connection.ChannelA.Port),
			fmt.Sprintf("recv_packet.packet_src_channel='%s'", connection.ChannelA.Channel),
			fmt.Sprintf("recv_packet.packet_dst_port='%s'", connection.ChannelB.Port),
			fmt.Sprintf("recv_packet.packet_dst_channel='%s'", connection.ChannelB.Channel),
		},
		"Receive IBC transaction",
		fmt.Sprintf("Receive IBC transaction for port: %s, channel: %s, sequence %s", connection.ChannelB.Port, connection.ChannelB.Channel, packetSequence),
		fmt.Sprintf("If the message does not get relayed, you can relay it yourself with hermes using the following command:\n hermes tx packet-recv --dst-chain %s --src-chain %s --src-port %s --src-channel %s\n", destinationChain.ChainID(), sourceChain.ChainID(), connection.ChannelA.Port, connection.ChannelA.Channel),
		timeoutHeight,
		timeoutTimestamp,
	)

	if timeout {
		fmt.Println("❌ IBC packet timed out")
		_, _ = waitForTXByEvents(
			cmd,
			sourceChain,
			[]string{
				fmt.Sprintf("timeout_packet.packet_sequence='%s'", packetSequence),
				fmt.Sprintf("timeout_packet.packet_src_port='%s'", connection.ChannelA.Port),
				fmt.Sprintf("timeout_packet.packet_src_channel='%s'", connection.ChannelA.Channel),
				fmt.Sprintf("timeout_packet.packet_dst_port='%s'", connection.ChannelB.Port),
				fmt.Sprintf("timeout_packet.packet_dst_channel='%s'", connection.ChannelB.Channel),
			},
			"Timeout/revert IBC transaction",
			"Timeout/revert IBC transaction",
			fmt.Sprintf("If the timeout message does not get relayed, you can relay it yourself with hermes using the following command:\n hermes tx packet-recv --dst-chain %s --src-chain %s --src-port %s --src-channel %s\n", destinationChain.ChainID(), sourceChain.ChainID(), connection.ChannelA.Port, connection.ChannelA.Channel),
			0,
			0,
		)
	}
}
