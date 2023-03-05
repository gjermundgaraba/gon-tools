package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"log"
)

func CreateClassCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-class",
		Short: "Create a NFT class",
		Long:  `Create a new class on any of the supported chains`,
		RunE: func(cmd *cobra.Command, args []string) error {
			wallet := chooseWallet(cmd)
			if err := cmd.Flags().Set(flags.FlagFrom, wallet); err != nil {
				panic(err)
			}

			chain := chooseChain("Select chain to create NFT on")
			setAddressPrefixes(chain.Bech32Prefix())

			clientCtx := getClientContext(cmd, chain)
			fromAccAddress := clientCtx.GetFromAddress()
			if fromAccAddress == nil {
				log.Fatal("No --from wallet/address specified")
			}
			fromAddress := chain.ConvertAccAddressToChainsPrefix(fromAccAddress)

			var classID string
			if err := survey.AskOne(&survey.Input{Message: "Class ID"}, &classID, survey.WithValidator(idValidator)); err != nil {
				return err
			}

			var symbol string
			if err := survey.AskOne(&survey.Input{Message: "Symbol"}, &symbol, survey.WithValidator(survey.Required)); err != nil {
				return err
			}

			var name string
			if err := survey.AskOne(&survey.Input{Message: "Name"}, &name, survey.WithValidator(survey.Required)); err != nil {
				return err
			}

			var description string
			if err := survey.AskOne(&survey.Input{Message: "Description"}, &description); err != nil {
				return err
			}

			var data string
			if err := survey.AskOne(&survey.Editor{
				Message:  "Data field (JSON)",
				FileName: "*.json",
			}, &data); err != nil {
				return err
			}

			var uri string
			if err := survey.AskOne(&survey.Input{Message: "URI"}, &uri); err != nil {
				return err
			}

			var uriHash string
			if uri != "" {
				if err := survey.AskOne(&survey.Input{Message: "URI Hash"}, &uriHash); err != nil {
					return err
				}
			}

			var mintRestricted bool
			var updateRestricted bool
			// TODO: find a less hardcoded and neat way to do this kind of stuff
			if chain.ChainID() == "gon-irishub-1" {
				if err := survey.AskOne(&survey.Confirm{Message: "Restrict minting NFTs to creator (you)?", Default: true}, &mintRestricted); err != nil {
					return err
				}

				if err := survey.AskOne(&survey.Confirm{Message: "Don't allow editing NFTs?", Default: true}, &updateRestricted); err != nil {
					return err
				}
			}

			msg := chain.CreateIssueCreditClassMsg(classID, name, "", fromAddress, symbol, mintRestricted, updateRestricted, description, uri, uriHash, data)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
