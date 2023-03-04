package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
)

type JunoChain struct {
	ChainData
}

func (c JunoChain) ListNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFT {
	panic("implement me")

}

func (c JunoChain) TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields) {
	panic("implement me")
}
