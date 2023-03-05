package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func CreateQueryClassesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classes",
		Short: "Query NFT classes",
		Long:  `Query for NFT classes on any of the supported chains`,
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

			classes := chain.ListNFTClasses(cmd.Context(), clientCtx, chains.ListNFTsQuery{
				Owner: fromAddress,
			})
			if len(classes) == 0 {
				log.Fatal("No NFT classes found")
			}

			for _, class := range classes {
				fmt.Printf("%s (%d nfts)\n", class.Label(), len(class.NFTs))
			}

			return nil
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
