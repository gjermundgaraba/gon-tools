package chains

import (
	"context"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

type omnfiFlixChain struct {
	ChainData
}

var OmniFlixChain = omnfiFlixChain{
	ChainData{
		name:              "OmniFlix GoN Testnet",
		chainID:           "gon-flixnet-1",
		bech32Prefix:      "omniflix",
		denom:             "uflix",
		rpc:               "http://65.21.93.56:26657",
		grpc:              "http://65.21.93.56:9090",
		nftImplementation: CosmosSDK,
	},
}

func (c omnfiFlixChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, _, _ bool, description, uri, uriHash, data string) sdk.Msg {
	return &omniflixnfttypes.MsgCreateDenom{
		Id:          denomID,
		Name:        denomName,
		Schema:      schema,
		Sender:      sender,
		Symbol:      symbol,
		Description: description,
		Uri:         uri,
		UriHash:     uriHash,
		Data:        data,
	}
}

func (c omnfiFlixChain) CreateTransferNFTMsg(connection NFTChannel, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
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

func (c omnfiFlixChain) CreateMintNFTMsg(tokenID, classID, tokenName, tokenURI, tokenURIHash, tokenData, minterAddress string) sdk.Msg {
	panic("implement me")
	/*return &omniflixnfttypes.MsgMintONFT{
		Id:           tokenID,
		DenomId:      classID,
		Metadata:     omniflixnfttypes.Metadata{},
		Data:         tokenData,
		Transferable: true,
		Extensible:   false,
		Nsfw:         false,
		RoyaltyShare: sdk.Dec{},
		Sender:       "",
		Recipient:    "",
	}*/
}

func (c omnfiFlixChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFTClass {
	nftQueryClient := omniflixnfttypes.NewQueryClient(clientCtx)

	request := &omniflixnfttypes.QueryOwnerONFTsRequest{
		DenomId: query.ClassReference,
		Owner:   query.Owner,
	}
	resp, err := nftQueryClient.OwnerONFTs(ctx, request)
	if err != nil {
		panic(err)
	}

	var classes []NFTClass
	for _, collection := range resp.Owner.IDCollections {
		var nfts []NFT
		for _, nft := range collection.OnftIds {
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
