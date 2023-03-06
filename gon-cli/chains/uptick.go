package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

type uptickChain struct {
	ChainData
}

var UptickChain = uptickChain{
	ChainData{
		name:         "Uptick GoN Testnet",
		chainID:      "uptick_7000-2",
		bech32Prefix: "uptick",
		denom:        "auptick",
		rpc:          "http://52.220.252.160:26657",
		grpc:         "http://52.220.252.160:9090",
	},
}

func (c uptickChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	panic("implement me")
}

func (c uptickChain) CreateTransferNFTMsg(connection NFTChannel, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	panic("implement me")
}

func (c uptickChain) CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg {
	panic("implement me")
}

func (c uptickChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFTClass {
	panic("implement me")
}
