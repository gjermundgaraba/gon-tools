package chains

import (
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
)

func createTransferNFTMsg(connection NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) sdk.Msg {
	return &nfttransfertypes.MsgTransfer{
		SourcePort:       connection.Port,
		SourceChannel:    connection.Channel,
		ClassId:          class.ClassID, // In the case of IBC, it will be the ibc/{hash} format
		TokenIds:         []string{nft.ID},
		Sender:           fromAddress,
		Receiver:         toAddress,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
		Memo:             "Sent using the Game of NFTs CLI by @gjermundgaraba",
	}
}
