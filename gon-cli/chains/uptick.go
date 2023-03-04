package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
)

type UptickChain struct {
	ChainData
}

func (c UptickChain) ListNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFT {
	panic("implement me")
}

func (c UptickChain) TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields) {
	panic("implement me")
}
