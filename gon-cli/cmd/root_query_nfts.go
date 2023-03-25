package cmd

import (
	"context"
	"fmt"
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/json"
	"log"
	"strings"
)

type queryNftContractResponse struct {
	Data string `json:"data"`
}

type queryNFTSOwnedResponse struct {
	Data struct {
		Tokens []string `json:"tokens"`
	} `json:"data"`
}

func queryNFTsInteractive(cmd *cobra.Command) error {
	chain := chooseChain("Select chain to create NFT on")
	setAddressPrefixes(chain.Bech32Prefix())

	key := chooseOrCreateKey(cmd, chain)
	if err := cmd.Flags().Set(flags.FlagFrom, key); err != nil {
		panic(err)
	}

	fromAddress := getAddressForChain(cmd, chain, key)

	clientCtx := getClientTxContext(cmd, chain)
	class := getUsersNfts(cmd.Context(), clientCtx, chain, fromAddress)
	fmt.Printf("Class ID: %s \n", class.ClassID)
	fmt.Printf("Base class ID: %s \n", class.BaseClassID)
	fmt.Printf("Full Path Class ID: %s \n", class.FullPathClassID)
	if chain.NFTImplementation() == chains.CosmWasm {
		fmt.Printf("Contract: %s \n", class.Contract)
	}
	if class.LastIBCChannel.Port != "" {
		fmt.Printf("Last IBC channel: %s/%s \n", class.LastIBCChannel.Port, class.LastIBCChannel.Channel)
	}

	fmt.Println()
	fmt.Println("NFTs:")
	for _, nft := range class.NFTs {
		fmt.Printf("- %s\n", nft.ID)
	}

	return nil
}

func getUsersNfts(ctx context.Context, clientCtx client.Context, chain chains.Chain, address string) chains.NFTClass {
	switch chain.NFTImplementation() {
	case chains.CosmosSDK:
		classes := chain.ListNFTClassesThatHasNFTs(ctx, clientCtx, address)
		if len(classes) == 0 {
			log.Fatal("No NFT classes found")
		}
		chosenClass := chooseOne("NFT Class", classes)

		return chosenClass
	case chains.CosmWasm:
		classID := askForString("Class ID (full IBC path/trace)")
		splitClassID := strings.Split(classID, "/")
		if len(splitClassID) <= 2 {
			panic("only IBC path class IDs are supported for CosmWasm chains right now")
		}

		lastPort := splitClassID[0]
		lastChannel := splitClassID[1]
		bridgerContract := strings.TrimPrefix(lastPort, "wasm.")

		nftContractQueryData, err := chains.Decoder.DecodeString(fmt.Sprintf(`{"nft_contract": {"class_id" : "%s"}}`, classID))
		if err != nil {
			panic(err)
		}
		wasmQueryClient := wasmdtypes.NewQueryClient(clientCtx)
		queryContractByClassIDResponse, err := wasmQueryClient.SmartContractState(
			ctx,
			&wasmdtypes.QuerySmartContractStateRequest{
				Address:   bridgerContract,
				QueryData: nftContractQueryData,
			},
		)
		if err != nil {
			panic(err)
		}
		queryContractStringOutput, err := clientCtx.Codec.MarshalJSON(queryContractByClassIDResponse)
		if err != nil {
			panic(err)
		}
		var nftContractResponse queryNftContractResponse
		if err := json.Unmarshal(queryContractStringOutput, &nftContractResponse); err != nil {
			panic(err)
		}
		nftContract := nftContractResponse.Data
		if nftContract == "" {
			log.Fatal("no NFT contract found, make sure you are using the full IBC path/trace with port/channel/class-id")
		}

		nftsOwnedQueryData, err := chains.Decoder.DecodeString(fmt.Sprintf(`{"tokens":{"owner": "%s"}}`, address))
		if err != nil {
			panic(err)
		}
		queryOwnedNFTsResponse, err := wasmQueryClient.SmartContractState(
			ctx,
			&wasmdtypes.QuerySmartContractStateRequest{
				Address:   nftContract,
				QueryData: nftsOwnedQueryData,
			},
		)
		if err != nil {
			panic(err)
		}
		queryNFTSStringOutput, err := clientCtx.Codec.MarshalJSON(queryOwnedNFTsResponse)
		if err != nil {
			panic(err)
		}
		var nftsOwnedResponse queryNFTSOwnedResponse
		if err := json.Unmarshal(queryNFTSStringOutput, &nftsOwnedResponse); err != nil {
			panic(err)
		}

		var nfts []chains.NFT

		for _, nft := range nftsOwnedResponse.Data.Tokens {
			nfts = append(nfts, chains.NFT{
				ID: nft,
			})
		}
		return chains.NFTClass{
			ClassID:         classID,
			BaseClassID:     splitClassID[len(splitClassID)-1],
			FullPathClassID: classID,
			Contract:        nftContract,
			NFTs:            nfts,
			LastIBCChannel: chains.NFTChannel{
				ChainID: chain.ChainID(),
				Port:    lastPort,
				Channel: lastChannel,
			},
		}
	default:
		panic("unknown NFT implementation")
	}
}
