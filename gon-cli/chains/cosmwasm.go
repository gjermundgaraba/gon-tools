package chains

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	"strings"
)

func createCosmWasmTransferMsg(connection NFTChannel, class NFTClass, nft NFT, fromAddress string, toAddress string, timeoutHeight clienttypes.Height) sdk.Msg {
	bridgerContract := strings.TrimPrefix(connection.Port, "wasm.")
	ibcTransferMsg := fmt.Sprintf(`{
  "receiver": "%s",
  "channel_id": "%s",
  "timeout": {
    "block": {
      "revision": %d,
      "height": %d
    }
  }
}`, toAddress, connection.Channel, timeoutHeight.RevisionNumber, timeoutHeight.RevisionHeight)
	ibcTransferMsgBase64Encoded := base64.StdEncoding.EncodeToString([]byte(ibcTransferMsg))

	execMsg := fmt.Sprintf(`{
  "send_nft": {
    "contract": "%s", 
    "token_id": "%s", 
    "msg": "%s"}
}`, bridgerContract, nft.ID, ibcTransferMsgBase64Encoded)
	return &wasmdtypes.MsgExecuteContract{
		Sender:   fromAddress,
		Contract: class.Contract,
		Msg:      []byte(execMsg),
	}
}

var Decoder = newArgDecoder(asciiDecodeString)

type argumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {
			return nil, fmt.Errorf("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return asciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func asciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
