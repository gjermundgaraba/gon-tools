package chains

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

type UptickChain struct {
	ChainData
}

func (c UptickChain) CreateIssueCreditClassMsg(denomID, denomName, schema, sender, symbol string, mintRestricted, updateRestricted bool, description, uri, uriHash, data string) sdk.Msg {
	panic("implement me")
}

func (c UptickChain) CreateTransferNFTMsg(connection NFTConnection, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	panic("implement me")
}

func (c UptickChain) ListNFTClasses(ctx context.Context, clientCtx client.Context, query ListNFTsQuery) []NFTClass {
	panic("implement me")
}
