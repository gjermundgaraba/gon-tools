package cmd

import (
	"fmt"
	"github.com/gjermundgaraba/goncli/chains"
	"github.com/spf13/cobra"
	"log"
)

func queryNFTClasses(cmd *cobra.Command) error {
	chain := chooseChain("Select chain to create NFT on")
	setAddressPrefixes(chain.Bech32Prefix())

	clientCtx := getClientContext(cmd, chain)
	fromAccAddress := clientCtx.GetFromAddress()
	if fromAccAddress == nil {
		log.Fatal("No --from wallet/address specified")
	}
	fromAddress := chain.ConvertAccAddressToChainsPrefix(fromAccAddress)

	classes := chain.ListNFTClassesThatHasNFTs(cmd.Context(), clientCtx, chains.ListNFTsQuery{
		Owner: fromAddress,
	})
	if len(classes) == 0 {
		log.Fatal("No NFT classes found")
	}

	for _, class := range classes {
		fmt.Printf("%s (%d nfts)\n", class.Label(), len(class.NFTs))
	}

	return nil
}
