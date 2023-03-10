package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	// Taken from IRIS: https://github.com/irisnet/irismod/blob/main/modules/nft/types/validation.go
	// DenomID or TokenID can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number or a separator ('/', ':', '.', '_' or '-').
	idString = `[a-z][a-zA-Z0-9/]{2,127}`
	regexpID = regexp.MustCompile(fmt.Sprintf(`^%s$`, idString)).MatchString
)

func idValidator(val interface{}) error {
	// since we are validating an Input, the assertion will always succeed
	if str, ok := val.(string); !ok || !regexpID(str) {
		return fmt.Errorf("ClassID can only accept characters that match the regular expression: %s", idString)
	}
	return nil
}

func setAddressPrefixes(prefix string) {
	accountPubKeyPrefix := prefix + "pub"
	validatorAddressPrefix := prefix + "valoper"
	validatorPubKeyPrefix := prefix + "valoperpub"
	consNodeAddressPrefix := prefix + "valcons"
	consNodePubKeyPrefix := prefix + "valconspub"

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(prefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
}

func getClientTxContext(cmd *cobra.Command, chain chains.Chain) client.Context {
	if err := cmd.Flags().Set(flags.FlagNode, chain.RPC()); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set(flags.FlagFees, fmt.Sprintf("25000%s", chain.Denom())); err != nil {
		panic(err)
	}
	if err := cmd.Flags().Set(flags.FlagGas, "1000000"); err != nil {
		panic(err)
	}

	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		panic(err)
	}

	return clientCtx.WithChainID(string(chain.ChainID()))
}

func getQueryClientContext(cmd *cobra.Command, chain chains.Chain) client.Context {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		panic(err)
	}

	newClient, err := client.NewClientFromNode(chain.RPC())
	if err != nil {
		panic(err)
	}

	return clientCtx.
		WithChainID(string(chain.ChainID())).
		WithNodeURI(chain.RPC()).
		WithClient(newClient)
}

func askForString(question string, opts ...survey.AskOpt) (answer string) {
	if err := survey.AskOne(&survey.Input{Message: question}, &answer, opts...); err != nil {
		panic(err)
	}

	return
}

func askForConfirmation(question string) bool {
	var answer bool
	if err := survey.AskOne(&survey.Confirm{Message: question}, &answer); err != nil {
		panic(err)
	}

	return answer
}

func chooseChain(questionPhrasing string, filterOut ...chains.Chain) chains.Chain {
	var chainOptions []chains.Chain
	for _, chain := range chains.Chains {
		toBeFilteredOut := false
		for _, filter := range filterOut {
			if chain.ChainID() == filter.ChainID() {
				toBeFilteredOut = true
				break
			}
		}
		if toBeFilteredOut {
			continue
		}

		chainOptions = append(chainOptions, chain)
	}

	return chooseOne(questionPhrasing, chainOptions)
}

func chooseConnection(sourceChain chains.Chain, destinationChain chains.Chain, chooseChannelQuestion string) chains.NFTConnection {
	connections := sourceChain.GetConnectionsTo(destinationChain)
	var wrappedConnections []OptionWrapper[chains.NFTConnection]
	for _, connection := range connections {
		wrappedConnections = append(wrappedConnections, OptionWrapper[chains.NFTConnection]{
			WrappedValue: connection,
			LabelFunc: func(connection chains.NFTConnection) string {
				return connection.ChannelA.Label()
			},
		})
	}

	return chooseOne(chooseChannelQuestion, wrappedConnections).WrappedValue
}

type Option interface {
	Label() string
}

type OptionString string

func (o OptionString) Label() string {
	return string(o)
}

type OptionWrapper[T any] struct {
	WrappedValue T
	LabelFunc    func(T) string
}

func (o OptionWrapper[T]) Label() string {
	return o.LabelFunc(o.WrappedValue)
}

func chooseOne[T Option](questionPhrasing string, options []T) T {
	var selectedIndex int
	var surveyOptions []string
	for _, o := range options {
		surveyOptions = append(surveyOptions, o.Label())
	}
	if err := survey.AskOne(&survey.Select{
		Message: questionPhrasing,
		Options: surveyOptions,
	}, &selectedIndex, survey.WithValidator(survey.Required)); err != nil {
		log.Fatalf("Error selecting: %v", err)
	}

	return options[selectedIndex]
}

// This is mostly copy-pasted from the Cosmos SDK, with the only difference being that it returns the tx response instead of printing it.
func sendTX(clientCtx client.Context, flagSet *pflag.FlagSet, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf := tx.NewFactoryCLI(clientCtx, flagSet)

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, err
		}
	}

	txf, err := txf.Prepare(clientCtx)
	if err != nil {
		return nil, err
	}

	if txf.SimulateAndExecute() || clientCtx.Simulate {
		_, adjusted, err := tx.CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", tx.GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if clientCtx.Simulate {
		return nil, nil
	}

	unsignedTx, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	if !clientCtx.SkipConfirm {
		txBytes, err := clientCtx.TxConfig.TxJSONEncoder()(unsignedTx.GetTx())
		if err != nil {
			return nil, err
		}

		if err := clientCtx.PrintRaw(txBytes); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", txBytes)
		}

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return nil, err
		}
	}

	err = tx.Sign(txf, clientCtx.GetFromName(), unsignedTx, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(unsignedTx.GetTx())
	if err != nil {
		return nil, err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func waitForTX(cmd *cobra.Command, chain chains.Chain, txHash string, txLabel string) *sdk.TxResponse {
	clientCtx := getQueryClientContext(cmd, chain)

	try := 1
	maxTries := 50
	for {
		if try > maxTries {
			panic(fmt.Errorf("%s (%s) on %s not found after %d tries", txLabel, txHash, chain.Label(), maxTries))
		}

		txResp, err := authtx.QueryTx(clientCtx, txHash)
		if err != nil {
			fmt.Print("\033[G\033[K") // move the cursor left and clear the line
			fmt.Printf("⬜ Waiting for %s (%s) on %s - attempt %d/%d", txLabel, txHash, chain.Label(), try, maxTries)
			time.Sleep(1 * time.Second)
			try++
			continue
		}

		if txResp.Code != 0 {
			panic(fmt.Errorf("transaction failed: %s", txResp.RawLog))
		}

		fmt.Print("\033[G\033[K") // move the cursor left and clear the line
		fmt.Printf("✅ %s (%s on %s) successful!\n", txLabel, chain.Label(), txHash)
		return txResp
	}
}

func waitForTXByEvents(cmd *cobra.Command, chain chains.Chain, events []string, txLabel string, timeoutHeight uint64, timeoutTimestamp uint64) (tx *sdk.TxResponse, timeout bool) {
	clientCtx := getQueryClientContext(cmd, chain)

	try := 1
	maxTries := 50
	for {
		if try > maxTries {
			panic(fmt.Errorf("%s not found after %d tries", txLabel, maxTries))
		}
		txsResult, err := authtx.QueryTxsByEvents(clientCtx, events, 1, 100, "asc")
		if err != nil {
			log.Fatalf("Error querying txs: %v", err)
		}

		switch len(txsResult.Txs) {
		case 0:
			if timeoutHeight != 0 || timeoutTimestamp != 0 {
				currentHeight, currentTimestamp := getCurrentChainStatus(cmd.Context(), clientCtx)
				if timeoutHeight != 0 && currentHeight >= timeoutHeight {
					return nil, true
				}
				if timeoutTimestamp != 0 && currentTimestamp >= timeoutTimestamp {
					return nil, true
				}
			}

			fmt.Print("\033[G\033[K") // move the cursor left and clear the line
			fmt.Printf("⬜ Waiting for %s on %s - attempt %d/%d", txLabel, chain.Label(), try, maxTries)
			time.Sleep(1 * time.Second)
			try++
			continue
		case 1:
			fmt.Print("\033[G\033[K") // move the cursor left and clear the line
			fmt.Printf("✅ %s (%s on %s) successful!\n", txLabel, txsResult.Txs[0].TxHash, chain.Label())
			return txsResult.Txs[0], false
		default:
			panic(fmt.Errorf("found more than one transaction with the same events"))
		}
	}
}

func getCurrentChainStatus(ctx context.Context, clientCtx client.Context) (height, timestamp uint64) {
	header, err := clientCtx.Client.Status(ctx)
	if err != nil {
		log.Fatalf("Error getting header: %v", err)
	}

	return uint64(header.SyncInfo.LatestBlockHeight), uint64(header.SyncInfo.LatestBlockTime.Nanosecond())
}

func calculateClassTrace(currentFullPathClassID string, connection chains.NFTConnection) (trace string, isRewind bool) {
	if strings.HasPrefix(currentFullPathClassID, fmt.Sprintf("%s/%s", connection.ChannelA.Port, connection.ChannelA.Channel)) {
		return strings.TrimPrefix(currentFullPathClassID, fmt.Sprintf("%s/%s/", connection.ChannelA.Port, connection.ChannelA.Channel)), true
	} else {
		return fmt.Sprintf("%s/%s/%s", connection.ChannelB.Port, connection.ChannelB.Channel, currentFullPathClassID), false
	}
}
