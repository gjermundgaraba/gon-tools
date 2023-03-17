package chains

import (
	"context"
	"log"
	"time"

	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channelutils "github.com/cosmos/ibc-go/v5/modules/core/04-channel/client/utils"
)

type NFTImplementation int

const (
	CosmosSDK NFTImplementation = iota
	CosmWasm
)

type KeyAlgo int

const (
	KeyAlgoSecp256k1 KeyAlgo = iota
	KeyAlgoEthSecp256k1
)

type Chain interface {
	Name() string
	Label() string
	ChainID() ChainID
	RPC() string
	GRPC() string
	Bech32Prefix() string
	Denom() string
	KeyAlgo() KeyAlgo
	NFTImplementation() NFTImplementation

	GetConnectionsTo(chain Chain) []NFTConnection
	GetIBCTimeouts(clientCtx client.Context, srcPort, srcChannel string, targetChainHeight, targetChainTimestamp uint64, tryToForceTimeout bool) (timeoutHeight clienttypes.Height, timeoutTimestamp uint64)

	CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg
	CreateTransferNFTMsg(channel NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg
	CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg

	ListNFTClassesThatHasNFTs(ctx context.Context, clientContext client.Context, owner string) []NFTClass
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
	ID string
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
	keyAlgo           KeyAlgo
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

func (c ChainData) KeyAlgo() KeyAlgo {
	return c.keyAlgo
}

func (c ChainData) NFTImplementation() NFTImplementation {
	return c.nftImplementation
}

func (c ChainData) GetIBCTimeouts(clientCtx client.Context, srcPort, srcChannel string, targetChainHeight, targetChainTimestamp uint64, tryToForceTimeout bool) (timeoutHeight clienttypes.Height, timeoutTimestamp uint64) {
	timeoutTimestamp = nfttransfertypes.DefaultRelativePacketTimeoutTimestamp
	timeoutHeight, err := clienttypes.ParseHeight(nfttransfertypes.DefaultRelativePacketTimeoutHeight)
	if err != nil {
		log.Fatalf("Error parsing timeout height: %v", err)
	}
	_, height, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
	if err != nil {
		log.Fatalf("Error querying latest consensus state: %v", err)
	}

	if tryToForceTimeout {
		return clienttypes.Height{
			RevisionNumber: height.RevisionNumber,
			RevisionHeight: height.RevisionHeight + 2,
		}, 0
	}

	absoluteHeight := height
	absoluteHeight.RevisionNumber += timeoutHeight.RevisionNumber
	absoluteHeight.RevisionHeight = targetChainHeight + timeoutHeight.RevisionHeight
	timeoutHeight = absoluteHeight

	// use local clock time as reference time if it is later than the
	// consensus state timestamp of the counter party chain, otherwise
	// still use consensus state timestamp as reference
	now := time.Now().UnixNano()
	if now <= 0 {
		log.Fatal("local clock time is not greater than Jan 1st, 1970 12:00 AM")
	}

	if uint64(now) > targetChainTimestamp {
		timeoutTimestamp = uint64(now) + timeoutTimestamp
	} else {
		timeoutTimestamp = targetChainTimestamp + timeoutTimestamp
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

func GetChainFromChainID(chainID ChainID) Chain {
	for _, chain := range Chains {
		if chain.ChainID() == chainID {
			return chain
		}
	}
	return nil
}
