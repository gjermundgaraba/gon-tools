package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
)

type IrisChain struct {
	ChainData
}

func (c IrisChain) ListNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFT {
	nftQueryClient := irisnfttypes.NewQueryClient(clientCtx)

	request := &irisnfttypes.QueryNFTsOfOwnerRequest{
		DenomId: query.ClassReference,
		Owner:   query.Owner,
	}
	resp, err := nftQueryClient.NFTsOfOwner(ctx, request)
	if err != nil {
		panic(err)
	}

	var nfts []NFT
	for _, collection := range resp.Owner.IDCollections {
		baseClassID, fullPathClassID, lastIBCConnection := findClassIBCInfo(ctx, clientCtx, collection.DenomId)

		for _, nft := range collection.TokenIds {
			nfts = append(nfts, NFT{
				ID:                nft,
				ClassID:           collection.DenomId,
				BaseClassID:       baseClassID,
				FullPathClassID:   fullPathClassID,
				LastIBCConnection: lastIBCConnection,
			})
		}
	}

	return nfts
}

func (c IrisChain) TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields) {
	panic("implement me")
}
