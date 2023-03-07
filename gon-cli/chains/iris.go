package chains

import (
	"context"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
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
		rpc:               "http://34.80.93.133:26657",
		grpc:              "http://34.80.93.133:9090",
		nftImplementation: CosmosSDK,
	},
}

func (c irisChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	return irisnfttypes.NewMsgIssueDenom(denomID, denomName, schema, sender, symbol, mintRestricted, updateRestricted, description, uri, uriHash, data)
}

func (c irisChain) CreateTransferNFTMsg(connection NFTChannel, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
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
				ID:      nft,
				ClassID: collection.DenomId,
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
