package chains

import "fmt"

type NFTChannel struct {
	ChainID ChainID
	Port    string
	Channel string
}

func (c NFTChannel) Label() string {
	return fmt.Sprintf("%s/%s", c.Port, c.Channel)
}

type NFTConnection struct {
	ChannelA NFTChannel
	ChannelB NFTChannel
}

// GetConnectionsTo returns all connections that have the given chain as the destination
// The output will be ordered so that ChannelA is always the one on the source chain
func (c ChainData) GetConnectionsTo(destinationChain Chain) []NFTConnection {
	var connections []NFTConnection
	for _, connection := range Connections {
		if connection.ChannelA.ChainID == c.ChainID() && connection.ChannelB.ChainID == destinationChain.ChainID() {
			connections = append(connections, connection)
		} else if connection.ChannelB.ChainID == c.ChainID() && connection.ChannelA.ChainID == destinationChain.ChainID() {
			// Switch them around so that the first channel is always the one on the source chain
			connections = append(connections, NFTConnection{
				ChannelA: connection.ChannelB,
				ChannelB: connection.ChannelA,
			})
		}
	}
	return connections
}

var Connections = []NFTConnection{
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-22",
		},
		ChannelB: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-207",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-23",
		},
		ChannelB: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-208",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-24",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-89",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-25",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-90",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-17",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-3",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-19",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-4",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-0",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-24",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: IRISChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-1",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-25",
		},
	},
	{

		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-230",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-120",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-234",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-122",
		},
	},
	{

		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-211",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-93",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-213",
		},
		ChannelB: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-94",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-203",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-6",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-206",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-12",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-209",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-44",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: StargazeChain.chainID,
			Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
			Channel: "channel-210",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-45",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-86",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-7",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-88",
		},
		ChannelB: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-13",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-91",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-46",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: JunoChain.chainID,
			Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
			Channel: "channel-92",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-47",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-5",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-41",
		},
	},
	{
		ChannelA: NFTChannel{
			ChainID: UptickChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-9",
		},
		ChannelB: NFTChannel{
			ChainID: OmniFlixChain.chainID,
			Port:    "nft-transfer",
			Channel: "channel-42",
		},
	},
}
