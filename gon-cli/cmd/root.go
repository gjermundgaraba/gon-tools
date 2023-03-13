package cmd

import (
	"fmt"
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctypes "github.com/cosmos/ibc-go/v5/modules/core/types"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/spf13/cobra"
	"github.com/strangelove-ventures/lens/client/codecs/ethermint"
	tmcfg "github.com/tendermint/tendermint/config"
	"os"
)

const (
	flagTryToForceTimeout = "try-to-timeout"

	createNFTClassOption  OptionString = "Create NFT Class"
	createNFTClassCommand              = "create-class"

	mintNFTOption  OptionString = "Mint NFT"
	mintNFTCommand              = "mint"

	queryNFTSOption  OptionString = "Query NFTs you own"
	queryNFTSCommand              = "query-nfts"

	transferNFTOption  OptionString = "Transfer NFT (Over IBC)"
	transferNFTCommand              = "transfer"

	toolsOption  OptionString = "Helper tools"
	toolsCommand              = "tools"

	gonToolsOption  OptionString = "GoN specific Tools"
	gonToolsCommand              = "gon-tools"
)

func NewRootCmd(appHomeDir string) *cobra.Command {
	encodingConfig := makeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(appHomeDir).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "gon [optional-command]",
		Short: "Game of NFTs - made simple!",
		Long: fmt.Sprintf(`Game of NFTs - made simple!
[optional-command] can be one of the following:
- %s (creates a new NFT class)
- %s (mints a new NFT)
- %s (queries your NFTs)
- %s (transfers an NFT over IBC)
- %s (lists available connections between to chains)
`, createNFTClassCommand, mintNFTCommand, queryNFTSCommand, transferNFTCommand, listConnectionsCommand),
		Args: cobra.ArbitraryArgs,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// set the default command outputs
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			// Overwrite here, because config.ReadFromClientConfig sets it...
			initClientCtx = initClientCtx.WithKeyring(getKeyring(encodingConfig.Codec))

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := tmcfg.DefaultConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			topLevelOptions := []OptionString{
				createNFTClassOption,
				mintNFTOption,
				transferNFTOption,
				queryNFTSOption,
				toolsOption,
				gonToolsOption,
			}

			var topLevelChoice OptionString
			if len(args) > 0 && args[0] != "" {
				switch args[0] {
				case createNFTClassCommand:
					topLevelChoice = createNFTClassOption
				case mintNFTCommand:
					topLevelChoice = mintNFTOption
				case transferNFTCommand:
					topLevelChoice = transferNFTOption
				case queryNFTSCommand:
					topLevelChoice = queryNFTSOption
				case toolsCommand:
					topLevelChoice = toolsOption
				case gonToolsCommand:
					topLevelChoice = gonToolsOption
				default:
					panic("invalid command")
				}
			} else {
				topLevelChoice = chooseOne("What would you like to do?", topLevelOptions)
			}

			switch topLevelChoice {
			case createNFTClassOption:
				return createNFTClassInteractive(cmd)
			case mintNFTOption:
				return mintNFTInteractive(cmd)
			case transferNFTOption:
				return transferNFTInteractive(cmd)
			case queryNFTSOption:
				return queryNFTsInteractive(cmd)
			case toolsOption:
				toolsInteractive(cmd, args)
				return nil
			case gonToolsOption:
				gonToolsInteractive(cmd, args)
				return nil
			default:
				panic(topLevelChoice + " not implemented option")
			}
			return nil
		},
	}

	flags.AddTxFlagsToCmd(rootCmd)
	rootCmd.Flags().Bool(flagTryToForceTimeout, false, "Try to force a timeout")

	return rootCmd
}

type EncodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Codec             codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// initAppConfig helps to override default appConfig template and configs.
// return "", nil if no custom configuration is required for the application.
func initAppConfig() (string, interface{}) {
	// The following code snippet is just for reference.
	type CustomAppConfig struct {
		serverconfig.Config
	}

	srvCfg := serverconfig.DefaultConfig()
	srvCfg.MinGasPrices = "0stake"

	customAppConfig := CustomAppConfig{
		Config: *srvCfg,
	}

	customAppTemplate := serverconfig.DefaultConfigTemplate

	return customAppTemplate, customAppConfig
}

func makeEncodingConfig() EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := tx.NewTxConfig(marshaler, tx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(amino)
	std.RegisterInterfaces(interfaceRegistry)

	ibctypes.RegisterInterfaces(interfaceRegistry)
	nfttransfertypes.RegisterInterfaces(interfaceRegistry)

	authtypes.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterLegacyAminoCodec(amino)

	irisnfttypes.RegisterInterfaces(interfaceRegistry)
	irisnfttypes.RegisterLegacyAminoCodec(amino)

	omniflixnfttypes.RegisterInterfaces(interfaceRegistry)
	omniflixnfttypes.RegisterLegacyAminoCodec(amino)

	wasmdtypes.RegisterInterfaces(interfaceRegistry)
	wasmdtypes.RegisterLegacyAminoCodec(amino)

	ethermint.RegisterInterfaces(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}
