package chains

import (
	"context"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channelutils "github.com/cosmos/ibc-go/v5/modules/core/04-channel/client/utils"
	"log"
	"time"
)

type Chain interface {
	Name() string
	ChainID() ChainID
	RPC() string
	GRPC() string
	Bech32Prefix() string
	Denom() string
	ConvertAddressToChainsPrefix(address string) string
	ConvertAccAddressToChainsPrefix(acc sdk.AccAddress) string
	GetIBCTimeouts(clientCtx client.Context, srcPort, srcChannel string) (timeoutHeight clienttypes.Height, timeoutTimestamp uint64)
	GetSourceNFTConnection(destinationChain Chain) NFTConnection
	ListNFTs(ctx context.Context, clientContext client.Context, query ListNFTsQuery) []NFT
	TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields)
}

type ListNFTsQuery struct {
	ClassReference string
	Owner          string
}

type TransferNFTFields struct {
	NFT              NFT
	DestinationChain Chain
	SenderAddress    string
	ReceiverAddress  string
}

type NFT struct {
	ID      string
	ClassID string
}

type ChainID string

type NFTType int

type NFTConnection struct {
	Port    string
	Channel string
}

type ChainData struct {
	name           string
	chainID        ChainID
	bech32Prefix   string
	denom          string
	rpc            string
	grpc           string
	nftType        NFTType
	nftConnections map[ChainID]NFTConnection
}

func (c ChainData) Name() string {
	return c.name
}

func (c ChainData) ChainID() ChainID {
	return c.chainID
}

func (c ChainData) RPC() string {
	return c.rpc
}

func (c ChainData) GRPC() string {
	return c.grpc
}

func (c ChainData) Bech32Prefix() string {
	return c.bech32Prefix
}

func (c ChainData) Denom() string {
	return c.denom
}

func (c ChainData) ConvertAddressToChainsPrefix(address string) string {
	_, acc, err := bech32.DecodeAndConvert(address)
	if err != nil {
		log.Fatalf("Error converting address: %v", err)
	}
	convertedAddress, err := bech32.ConvertAndEncode(c.bech32Prefix, acc)
	if err != nil {
		log.Fatalf("Error converting address: %v", err)
	}

	return convertedAddress
}

func (c ChainData) ConvertAccAddressToChainsPrefix(acc sdk.AccAddress) string {
	convertedAddress, err := bech32.ConvertAndEncode(c.bech32Prefix, acc.Bytes())
	if err != nil {
		log.Fatalf("Error converting address: %v", err)
	}

	return convertedAddress
}

func (c ChainData) GetIBCTimeouts(clientCtx client.Context, srcPort, srcChannel string) (timeoutHeight clienttypes.Height, timeoutTimestamp uint64) {
	timeoutTimestamp = nfttransfertypes.DefaultRelativePacketTimeoutTimestamp
	timeoutHeight, err := clienttypes.ParseHeight(nfttransfertypes.DefaultRelativePacketTimeoutHeight)
	if err != nil {
		log.Fatalf("Error parsing timeout height: %v", err)
	}

	consensusState, height, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
	if err != nil {
		log.Fatalf("Error querying latest consensus state: %v", err)
	}

	absoluteHeight := height
	absoluteHeight.RevisionNumber += timeoutHeight.RevisionNumber
	absoluteHeight.RevisionHeight += timeoutHeight.RevisionHeight
	timeoutHeight = absoluteHeight

	// use local clock time as reference time if it is later than the
	// consensus state timestamp of the counter party chain, otherwise
	// still use consensus state timestamp as reference
	now := time.Now().UnixNano()
	if now <= 0 {
		log.Fatal("local clock time is not greater than Jan 1st, 1970 12:00 AM")
	}

	consensusStateTimestamp := consensusState.GetTimestamp()
	if uint64(now) > consensusStateTimestamp {
		timeoutTimestamp = uint64(now) + timeoutTimestamp
	} else {
		timeoutTimestamp = consensusStateTimestamp + timeoutTimestamp
	}

	return
}

func (c ChainData) GetSourceNFTConnection(destinationChain Chain) NFTConnection {
	return c.nftConnections[destinationChain.ChainID()]
}

const (
	NFTTypeCosmosSDK NFTType = iota
	NFTTypeWasm
)

var Chains = []Chain{
	IrisChain{
		ChainData{
			name:         "IRISNet GoN Testnet",
			chainID:      "gon-irishub-1",
			bech32Prefix: "iaa",
			denom:        "uiris",
			rpc:          "http://34.80.93.133:26657",
			grpc:         "http://34.80.93.133:9090",
			nftType:      NFTTypeCosmosSDK,
			nftConnections: map[ChainID]NFTConnection{
				"elgafar-1": {
					Port:    "nft-transfer",
					Channel: "channel-22",
				},
				"uni-6": {
					Port:    "nft-transfer",
					Channel: "channel-24",
				},
				"uptick_7000-2": {
					Port:    "nft-transfer",
					Channel: "channel-17",
				},
				"gon-flixnet-1": {
					Port:    "nft-transfer",
					Channel: "channel-0",
				},
			},
		},
	},
	StargazeChain{
		ChainData{
			name:         "Stargaze GoN Testnet",
			chainID:      "elgafar-1",
			bech32Prefix: "stars",
			denom:        "ustars",
			rpc:          "https://rpc.elgafar-1.stargaze-apis.com:443",
			grpc:         "http://grpc-1.elgafar-1.stargaze-apis.com:26660",
			nftType:      NFTTypeWasm,
			nftConnections: map[ChainID]NFTConnection{
				"gon-irishub-1": {
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-207",
				},
				"uni-6": {
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-211",
				},
				"uptick_7000-2": {
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-203",
				},
				"gon-flixnet-1": {
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-209",
				},
			},
		},
	},
	JunoChain{
		ChainData{
			name:         "Juno GoN Testnet",
			chainID:      "uni-6",
			bech32Prefix: "juno",
			denom:        "ujunox",
			rpc:          "https://rpc.uni.junonetwork.io:443",
			grpc:         "http://juno-testnet-grpc.polkachu.com:12690",
			nftType:      NFTTypeWasm,
			nftConnections: map[ChainID]NFTConnection{
				"gon-irishub-1": {
					Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
					Channel: "channel-89",
				},
				"elgafar-1": {
					Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
					Channel: "channel-93",
				},
				"uptick_7000-2": {
					Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
					Channel: "channel-86",
				},
				"gon-flixnet-1": {
					Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
					Channel: "channel-91",
				},
			},
		},
	},
	UptickChain{
		ChainData{
			name:         "Uptick GoN Testnet",
			chainID:      "uptick_7000-2",
			bech32Prefix: "uptick",
			denom:        "auptick",
			rpc:          "http://52.220.252.160:26657",
			grpc:         "http://52.220.252.160:9090",
			nftType:      NFTTypeCosmosSDK,
			nftConnections: map[ChainID]NFTConnection{
				"gon-irishub-1": {
					Port:    "nft-transfer",
					Channel: "channel-3",
				},
				"elgafar-1": {
					Port:    "nft-transfer",
					Channel: "channel-6",
				},
				"uni-6": {
					Port:    "nft-transfer",
					Channel: "channel-7",
				},
				"gon-flixnet-1": {
					Port:    "nft-transfer",
					Channel: "channel-5",
				},
			},
		},
	},
	OmnfiFlixChain{
		ChainData{
			name:         "OmniFlix GoN Testnet",
			chainID:      "gon-flixnet-1",
			bech32Prefix: "omniflix",
			denom:        "uflix",
			rpc:          "http://65.21.93.56:26657",
			grpc:         "http://65.21.93.56:9090",
			nftType:      NFTTypeCosmosSDK,
			nftConnections: map[ChainID]NFTConnection{
				"gon-irishub-1": {
					Port:    "nft-transfer",
					Channel: "channel-24",
				},
				"elgafar-1": {
					Port:    "nft-transfer",
					Channel: "channel-44",
				},
				"uni-6": {
					Port:    "nft-transfer",
					Channel: "channel-46",
				},
				"uptick_7000-2": {
					Port:    "nft-transfer",
					Channel: "channel-41",
				},
			},
		},
	},
}
