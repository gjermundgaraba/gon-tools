package chains

import (
	"context"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
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
		keyAlgo:           KeyAlgoSecp256k1,
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

func (c omnfiFlixChain) CreateTransferNFTMsg(connection NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	return createTransferNFTMsg(connection, class, nft, fromAddress, toAddress, timeoutHeight, timeoutTimestamp)
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

func (c omnfiFlixChain) ListNFTClassesThatHasNFTs(ctx context.Context, clientCtx client.Context, owner string) []NFTClass {
	nftQueryClient := omniflixnfttypes.NewQueryClient(clientCtx)

	request := &omniflixnfttypes.QueryOwnerONFTsRequest{
		Owner: owner,
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
