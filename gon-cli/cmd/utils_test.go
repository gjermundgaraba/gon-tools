package cmd

import (
	"context"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"os"
	"path/filepath"
	"testing"
)

func TestSetAddressPrefix(t *testing.T) {
	accAddr, err := sdk.AccAddressFromBech32("cosmos1ypynejafjw6u2cucqp9yjnahxjt3vkl6wj2c72")
	require.NoError(t, err)

	setAddressPrefixes(chains.IRISChain.Bech32Prefix())
	require.Equal(t, chains.IRISChain.Bech32Prefix(), sdk.GetConfig().GetBech32AccountAddrPrefix())
	require.Equal(t, "iaa1ypynejafjw6u2cucqp9yjnahxjt3vkl6ms2fum", accAddr.String())

	setAddressPrefixes(chains.StargazeChain.Bech32Prefix())
	require.Equal(t, chains.StargazeChain.Bech32Prefix(), sdk.GetConfig().GetBech32AccountAddrPrefix())
	require.Equal(t, "stars1ypynejafjw6u2cucqp9yjnahxjt3vkl66wa94m", accAddr.String())
}

func TestNftClassFromTrace(t *testing.T) {
	testTable := map[string]struct {
		FullPathClassID  string
		DestinationChain chains.Chain
		ExpectedNFTClass chains.NFTClass
	}{
		"baseClassID": {
			FullPathClassID:  "baseclassid",
			DestinationChain: chains.OmniFlixChain,
			ExpectedNFTClass: chains.NFTClass{
				ClassID:         "baseclassid",
				BaseClassID:     "baseclassid",
				FullPathClassID: "baseclassid",
				Contract:        "",
				NFTs:            nil,
				LastIBCChannel:  chains.NFTChannel{},
			},
		},
		"iris to stargaze 1": {
			FullPathClassID:  "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-207/bugtest01",
			DestinationChain: chains.StargazeChain,
			ExpectedNFTClass: chains.NFTClass{
				ClassID:         "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-207/bugtest01",
				BaseClassID:     "bugtest01",
				FullPathClassID: "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-207/bugtest01",
				Contract:        "stars1hflxm65g37tshrmheg5ue5ufdj84zuvp22jsuq08q09pqwrv78fqpampz2",
				NFTs:            nil,
				LastIBCChannel: chains.NFTChannel{
					ChainID: "",
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-207",
				},
			},
		},
		"iris to stargaze 2": {
			FullPathClassID:  "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
			DestinationChain: chains.StargazeChain,
			ExpectedNFTClass: chains.NFTClass{
				ClassID:         "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
				BaseClassID:     "bugtest01",
				FullPathClassID: "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
				Contract:        "stars1sxj7p6vrsefv7s5j0acqfw37j27slpu6mtvnwx9mt7u5mauk74ysn45unj",
				NFTs:            nil,
				LastIBCChannel: chains.NFTChannel{
					ChainID: "",
					Port:    "wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh",
					Channel: "channel-208",
				},
			},
		},
		"stargaze to juno": {
			FullPathClassID:  "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120/wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
			DestinationChain: chains.JunoChain,
			ExpectedNFTClass: chains.NFTClass{
				ClassID:         "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120/wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
				BaseClassID:     "bugtest01",
				FullPathClassID: "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120/wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
				Contract:        "juno1zkhxyq5s63k0p40he3jjvvva0cdmac5v7pmehynqrhx2zaat9l7sdvn0gy",
				NFTs:            nil,
				LastIBCChannel: chains.NFTChannel{
					ChainID: "",
					Port:    "wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a",
					Channel: "channel-120",
				},
			},
		},
		"juno to uptick": {
			FullPathClassID:  "nft-transfer/channel-7/wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120/wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
			DestinationChain: chains.UptickChain,
			ExpectedNFTClass: chains.NFTClass{
				ClassID:         "ibc/F96B818A4B41C89AC456A86413FC306D4B56A8F9B3AE347A60F499A34AEC1B27",
				BaseClassID:     "bugtest01",
				FullPathClassID: "nft-transfer/channel-7/wasm.juno1stv6sk0mvku34fj2mqrlyru6683866n306mfv52tlugtl322zmks26kg7a/channel-120/wasm.stars1ve46fjrhcrum94c7d8yc2wsdz8cpuw73503e8qn9r44spr6dw0lsvmvtqh/channel-208/bugtest01",
				Contract:        "",
				NFTs:            nil,
				LastIBCChannel: chains.NFTChannel{
					ChainID: "",
					Port:    "nft-transfer",
					Channel: "channel-7",
				},
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			cmd := getTestCmd(t)

			actualNFTClass := queryNftClassFromTrace(cmd, test.FullPathClassID, test.DestinationChain)
			require.Equal(t, test.ExpectedNFTClass, actualNFTClass)
		})
	}
}

func getTestCmd(t *testing.T) *cobra.Command {
	userHomeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	appHomeDir := filepath.Join(userHomeDir, ".gon-cli")

	cmd := &cobra.Command{
		Use: "test",
		Run: func(_ *cobra.Command, _ []string) {},
	}
	flags.AddTxFlagsToCmd(cmd)
	encodingConfig := makeEncodingConfig()
	initClientCtx := getInitialClientCtx(appHomeDir)

	srvCtx := server.NewDefaultContext()
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})
	ctx = context.WithValue(ctx, server.ServerContextKey, srvCtx)
	cmd.SetContext(ctx)

	_ = tmcli.PrepareBaseCmd(cmd, "", appHomeDir)

	require.NoError(t, persistentPreRun(cmd, initClientCtx, encodingConfig.Codec))
	return cmd
}
