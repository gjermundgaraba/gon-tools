package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
)

type StargazeChain struct {
	ChainData
}

func (c StargazeChain) ListNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFT {
	panic("implement me")
}

func (c StargazeChain) TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields) {
	panic("implement me")
}
