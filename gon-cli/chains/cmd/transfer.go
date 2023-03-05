package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer NFTs over IBC",
		Long:  `Transfer NFTs from one address to another - over IBC`,
		RunE: func(cmd *cobra.Command, args []string) error {
			sourceChain := chooseChain("Select source chain")
			setAddressPrefixes(sourceChain.Bech32Prefix())

			clientCtx := getClientContext(cmd, sourceChain)
			fromAccAddress := clientCtx.GetFromAddress()
			if fromAccAddress == nil {
				log.Fatal("No --from wallet/address specified")
			}

			fromAddress := sourceChain.ConvertAccAddressToChainsPrefix(fromAccAddress)

			destinationChain := chooseChain("Select destination chain", sourceChain)
			_ = destinationChain

			/*var classId string
			if err := survey.AskOne(&survey.Input{Message: "What is the Class ID on the source chain?"}, &classId, survey.WithValidator(survey.Required)); err != nil {
				log.Fatalf("Error getting class ID: %v", err)
			}*/

			nfts := sourceChain.ListNFTs(cmd.Context(), clientCtx, chains.ListNFTsQuery{
				Owner: fromAddress,
			})
			if len(nfts) == 0 {
				log.Fatal("No NFTs found")
			}

			classMap := make(map[string][]chains.NFT)
			for _, nft := range nfts {
				classMap[nft.FullPathClassID] = append(classMap[nft.FullPathClassID], nft)
			}

			var classId string
			var classOptions []string
			for classId := range classMap {
				classOptions = append(classOptions, classId)
			}
			if err := survey.AskOne(&survey.Select{
				Message: "Select class",
				Options: classOptions,
			}, &classId, survey.WithValidator(survey.Required)); err != nil {
				log.Fatalf("Error selecting class: %v", err)
			}

			// select nft
			var nftId string
			var nftOptions []string
			for _, nft := range classMap[classId] {
				nftOptions = append(nftOptions, nft.ID)
			}
			if err := survey.AskOne(&survey.Select{
				Message: "Select NFT",
				Options: nftOptions,
			}, &nftId, survey.WithValidator(survey.Required)); err != nil {
				log.Fatalf("Error selecting NFT: %v", err)
			}
			var selectedNFT chains.NFT
			for _, nft := range classMap[classId] {
				if nft.ID == nftId {
					selectedNFT = nft
					break
				}
			}

			var destinationAddress string
			if err := survey.AskOne(&survey.Input{Message: "What is the destination address? (Leave empty to send to same address as owner on destination)"}, &destinationAddress); err != nil {
				log.Fatalf("Error getting destination address: %v", err)
			}
			if destinationAddress == "" {
				destinationAddress = destinationChain.ConvertAddressToChainsPrefix(fromAddress)
			}

			/*sourceChain.TransferNFT(cmd.Context(), clientCtx, chains.TransferNFTFields{
				NFT: chains.NFT{
					ID:      nftId,
					FullPathClassID: classId,
				},
				DestinationChain: destinationChain,
				SenderAddress:    fromAddress,
				ReceiverAddress:  destinationAddress,
			})*/

			sourceNFTConnection := selectedNFT.LastIBCConnection
			if sourceNFTConnection.Port == "" {
				sourceNFTConnection = sourceChain.GetSourceNFTConnection(destinationChain)
			}
			timeoutHeight, timeoutTimestamp := sourceChain.GetIBCTimeouts(clientCtx, sourceNFTConnection.Port, sourceNFTConnection.Channel)

			msg := &nfttransfertypes.MsgTransfer{
				SourcePort:       sourceNFTConnection.Port,
				SourceChannel:    sourceNFTConnection.Channel,
				ClassId:          selectedNFT.ClassID, // In the case of IBC, it will be the ibc/{hash} format
				TokenIds:         []string{nftId},
				Sender:           fromAddress,
				Receiver:         destinationAddress,
				TimeoutHeight:    timeoutHeight,
				TimeoutTimestamp: timeoutTimestamp,
				Memo:             "Sent using the Game of NFTs CLI by gjermundgaraba",
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
