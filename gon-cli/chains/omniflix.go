package chains

import (
	"context"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
	"github.com/cosmos/cosmos-sdk/client"
)

type OmnfiFlixChain struct {
	ChainData
}

func (c OmnfiFlixChain) ListNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFT {
	nftQueryClient := omniflixnfttypes.NewQueryClient(clientCtx)

	request := &omniflixnfttypes.QueryOwnerONFTsRequest{
		DenomId: query.ClassReference,
		Owner:   query.Owner,
	}
	resp, err := nftQueryClient.OwnerONFTs(ctx, request)
	if err != nil {
		panic(err)
	}

	var nfts []NFT
	for _, collection := range resp.Owner.IDCollections {
		baseClassID, fullPathClassID, lastIBCConnection := findClassIBCInfo(ctx, clientCtx, collection.DenomId)

		for _, nft := range collection.OnftIds {
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

func (c OmnfiFlixChain) TransferNFT(ctx context.Context, clientCtx client.Context, fields TransferNFTFields) {
	panic("implement me")
}
