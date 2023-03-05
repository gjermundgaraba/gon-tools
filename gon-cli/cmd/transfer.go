package cmd

import (
	"github.com/AlecAivazis/survey/v2"
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
			wallet := chooseWallet(cmd)
			if err := cmd.Flags().Set(flags.FlagFrom, wallet); err != nil {
				panic(err)
			}
			
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

			classes := sourceChain.ListNFTClasses(cmd.Context(), clientCtx, chains.ListNFTsQuery{
				Owner: fromAddress,
			})
			if len(classes) == 0 {
				log.Fatal("No NFT classes found")
			}

			selectedClass := selectOne("Select class", classes)

			// select nft
			// TODO: Use multiselect to be able to send more than one at a time
			selectedNFT := selectOne("Select NFT", selectedClass.NFTs)

			var destinationAddress string
			if err := survey.AskOne(&survey.Input{Message: "What is the destination address? (Leave empty to send to same address as owner on destination)"}, &destinationAddress); err != nil {
				log.Fatalf("Error getting destination address: %v", err)
			}
			if destinationAddress == "" {
				destinationAddress = destinationChain.ConvertAddressToChainsPrefix(fromAddress)
			}

			sourceNFTConnection := selectedClass.LastIBCConnection
			if sourceNFTConnection.Port == "" {
				sourceNFTConnection = sourceChain.GetSourceNFTConnection(destinationChain)
			}
			timeoutHeight, timeoutTimestamp := sourceChain.GetIBCTimeouts(clientCtx, sourceNFTConnection.Port, sourceNFTConnection.Channel)

			msg := sourceChain.CreateTransferNFTMsg(sourceNFTConnection, selectedNFT, fromAddress, destinationAddress, timeoutHeight, timeoutTimestamp)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
