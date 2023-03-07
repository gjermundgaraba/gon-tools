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

type NFTImplementation int

const (
	CosmosSDK NFTImplementation = iota
	CosmWasm
)

type Chain interface {
	Name() string
	Label() string
	ChainID() ChainID
	RPC() string
	GRPC() string
	Bech32Prefix() string
	Denom() string
	NFTImplementation() NFTImplementation
	ConvertAddressToChainsPrefix(address string) string
	ConvertAccAddressToChainsPrefix(acc sdk.AccAddress) string

	GetConnectionsTo(chain Chain) []NFTConnection
	GetIBCTimeouts(clientCtx client.Context, srcPort, srcChannel string) (timeoutHeight clienttypes.Height, timeoutTimestamp uint64)

	CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg
	CreateTransferNFTMsg(channel NFTChannel, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg
	CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg

	ListNFTClassesThatHasNFTs(ctx context.Context, clientContext client.Context, query ListNFTsQuery) []NFTClass
}

type ListNFTsQuery struct {
	ClassReference string
	Owner          string
}

type NFTClass struct {
	ClassID         string
	BaseClassID     string
	FullPathClassID string
	Contract        string // CosmWasm only
	NFTs            []NFT
	LastIBCChannel  NFTChannel
}

func (n NFTClass) Label() string {
	return n.FullPathClassID
}

type NFT struct {
	ID      string
	ClassID string
}

func (n NFT) Label() string {
	return n.ID
}

type ChainID string

type ChainData struct {
	name              string
	chainID           ChainID
	bech32Prefix      string
	denom             string
	rpc               string
	grpc              string
	nftImplementation NFTImplementation
}

func (c ChainData) Name() string {
	return c.name
}

func (c ChainData) Label() string {
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

func (c ChainData) NFTImplementation() NFTImplementation {
	return c.nftImplementation
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

var Chains = []Chain{
	IRISChain,
	StargazeChain,
	JunoChain,
	UptickChain,
	OmniFlixChain,
}
