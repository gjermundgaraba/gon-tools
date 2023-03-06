package chains

import (
	"context"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
)

type IrisChain struct {
	ChainData
}

func (c IrisChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	return irisnfttypes.NewMsgIssueDenom(denomID, denomName, schema, sender, symbol, mintRestricted, updateRestricted, description, uri, uriHash, data)
}

func (c IrisChain) CreateTransferNFTMsg(connection NFTConnection, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	return &nfttransfertypes.MsgTransfer{
		SourcePort:       connection.Port,
		SourceChannel:    connection.Channel,
		ClassId:          nft.ClassID, // In the case of IBC, it will be the ibc/{hash} format
		TokenIds:         []string{nft.ID},
		Sender:           fromAddress,
		Receiver:         toAddress,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             "Sent using the Game of NFTs CLI by @gjermundgaraba",
	}
}

func (c IrisChain) ListNFTClasses(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFTClass {
	nftQueryClient := irisnfttypes.NewQueryClient(clientCtx)

	request := &irisnfttypes.QueryNFTsOfOwnerRequest{
		DenomId: query.ClassReference,
		Owner:   query.Owner,
	}
	resp, err := nftQueryClient.NFTsOfOwner(ctx, request)
	if err != nil {
		panic(err)
	}

	var classes []NFTClass
	for _, collection := range resp.Owner.IDCollections {
		var nfts []NFT
		for _, nft := range collection.TokenIds {
			nfts = append(nfts, NFT{
				ID:      nft,
				ClassID: collection.DenomId,
			})
		}

		baseClassID, fullPathClassID, lastIBCConnection := findClassIBCInfo(ctx, clientCtx, collection.DenomId)
		classes = append(classes, NFTClass{
			ClassID:           collection.DenomId,
			BaseClassID:       baseClassID,
			FullPathClassID:   fullPathClassID,
			NFTs:              nfts,
			LastIBCConnection: lastIBCConnection,
		})
	}

	return classes
}
