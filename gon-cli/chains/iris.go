package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
)

type irisChain struct {
	ChainData
}

var IRISChain = irisChain{
	ChainData{
		name:              "IRISNet GoN Testnet",
		chainID:           "gon-irishub-1",
		bech32Prefix:      "iaa",
		denom:             "uiris",
		keyAlgo:           KeyAlgoSecp256k1,
		rpc:               "http://34.80.93.133:26657",
		grpc:              "http://34.80.93.133:9090",
		nftImplementation: CosmosSDK,
	},
}

func (c irisChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	return irisnfttypes.NewMsgIssueDenom(denomID, denomName, schema, sender, symbol, mintRestricted, updateRestricted, description, uri, uriHash, data)
}

func (c irisChain) CreateTransferNFTMsg(connection NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	return createTransferNFTMsg(connection, class, nft, fromAddress, toAddress, timeoutHeight, timeoutTimestamp)
}

func (c irisChain) CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg {
	return irisnfttypes.NewMsgMintNFT(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress, minterAddress)
}

func (c irisChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, owner string) []NFTClass {
	nftQueryClient := irisnfttypes.NewQueryClient(clientCtx)

	request := &irisnfttypes.QueryNFTsOfOwnerRequest{
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
