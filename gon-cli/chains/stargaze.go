package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

type stargazeChain struct {
	ChainData
}

var StargazeChain = stargazeChain{
	ChainData{
		name:         "Stargaze GoN Testnet",
		chainID:      "elgafar-1",
		bech32Prefix: "stars",
		denom:        "ustars",
		rpc:          "https://rpc.elgafar-1.stargaze-apis.com:443",
		grpc:         "http://grpc-1.elgafar-1.stargaze-apis.com:26660",
	},
}

func (c stargazeChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	panic("implement me")
}

func (c stargazeChain) CreateTransferNFTMsg(connection NFTChannel, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	panic("implement me")
}

func (c stargazeChain) CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg {
	panic("implement me")
}

func (c stargazeChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFTClass {
	panic("implement me")
}
