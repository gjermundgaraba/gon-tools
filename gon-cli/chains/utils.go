package chains

import (
	"context"
	"fmt"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	"strings"
)

func findClassIBCInfo(ctx context.Context, clientCtx client.Context, classID string) (baseClassID string, fullPathClassID string, lastIBCConnection NFTConnection) {
	baseClassID = classID
	fullPathClassID = classID
	if strings.HasPrefix(classID, "ibc/") {
		classHash := strings.Split(classID, "/")[1]
		traceReq := &nfttransfertypes.QueryClassTraceRequest{
			Hash: classHash,
		}

		nftTransferQueryClient := nfttransfertypes.NewQueryClient(clientCtx)
		traceResp, err := nftTransferQueryClient.ClassTrace(ctx, traceReq)
		if err != nil {
			panic(err)
		}

		baseClassID = traceResp.ClassTrace.BaseClassId
		fullPathClassID = fmt.Sprintf("%s/%s", traceResp.ClassTrace.Path, baseClassID)
		pathSplit := strings.Split(traceResp.ClassTrace.Path, "/")
		latestPort := pathSplit[len(pathSplit)-2]
		latestChannel := pathSplit[len(pathSplit)-1]
		lastIBCConnection = NFTConnection{
			Port:    latestPort,
			Channel: latestChannel,
		}
	}

	return
}
