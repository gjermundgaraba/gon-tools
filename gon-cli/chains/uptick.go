package chains

import (
	"context"
	uptickcollectiontypes "github.com/UptickNetwork/uptick/x/collection/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

type uptickChain struct {
	ChainData
}

var UptickChain = uptickChain{
	ChainData{
		name:              "Uptick GoN Testnet",
		chainID:           "uptick_7000-2",
		bech32Prefix:      "uptick",
		denom:             "auptick",
		keyAlgo:           KeyAlgoEthSecp256k1,
		rpc:               "http://52.220.252.160:26657",
		grpc:              "http://52.220.252.160:9090",
		nftImplementation: CosmosSDK,
	},
}

func (c uptickChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	panic("implement me")
}

func (c uptickChain) CreateTransferNFTMsg(connection NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	return createTransferNFTMsg(connection, class, nft, fromAddress, toAddress, timeoutHeight, timeoutTimestamp)
}

func (c uptickChain) CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg {
	panic("implement me")
}

func (c uptickChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, owner string) []NFTClass {
	nftQueryClient := uptickcollectiontypes.NewQueryClient(clientCtx)

	request := &uptickcollectiontypes.QueryNFTsOfOwnerRequest{
		Owner: owner,
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
				ID: nft,
			})
		}

		baseClassID, fullPathClassID, lastIBCConnection := findClassIBCInfo(ctx, clientCtx, collection.DenomId)
		classes = append(classes, NFTClass{
			ClassID:         collection.DenomId,
			BaseClassID:     baseClassID,
			FullPathClassID: fullPathClassID,
			NFTs:            nfts,
			LastIBCChannel:  lastIBCConnection,
		})
	}

	return classes
}
