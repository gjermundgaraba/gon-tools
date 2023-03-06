package cmd

import (
	"fmt"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
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
	tmcfg "github.com/tendermint/tendermint/config"
	"os"
)

const (
	createNFTClassOption  OptionString = "Create NFT Class"
	createNFTClassCommand              = "create-class"

	mintNFTOption  OptionString = "Mint NFT"
	mintNFTCommand              = "mint"

	queryNFTClassesOption  OptionString = "Query NFT Classes for which you own NFTs"
	queryNFTClassesCommand              = "query-classes"

	transferNFTOption  OptionString = "Transfer NFT (Over IBC)"
	transferNFTCommand              = "transfer"
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
		Use:   "gon-cli [optional-command]",
		Short: "Game of NFTs - made simple!",
		Long: fmt.Sprintf(`Game of NFTs - made simple!
[optional-command] can be one of the following:
- %s (creates a new NFT class)
- %s (mints a new NFT)
- %s (queries your NFT classes)
- %s (transfers an NFT over IBC)
`, createNFTClassCommand, mintNFTCommand, queryNFTClassesCommand, transferNFTCommand),
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

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := initAppConfig()
			customTMConfig := tmcfg.DefaultConfig()

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig, customTMConfig)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			wallet := chooseWallet(cmd)
			if err := cmd.Flags().Set(flags.FlagFrom, wallet); err != nil {
				panic(err)
			}

			topLevelOptions := []OptionString{
				createNFTClassOption,
				mintNFTOption,
				transferNFTOption,
				queryNFTClassesOption,
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
				case queryNFTClassesCommand:
					topLevelChoice = queryNFTClassesOption
				default:
					panic("invalid command")
				}
			} else {
				topLevelChoice = chooseOne("What would you like to do?", topLevelOptions)
			}

			switch topLevelChoice {
			case createNFTClassOption:
				return createNFTClass(cmd)
			case mintNFTOption:
				return mintNFT(cmd)
			case transferNFTOption:
				return transferNFT(cmd)
			case queryNFTClassesOption:
				return queryNFTClasses(cmd)
			}
			return nil
		},
	}

	rootCmd.AddCommand(
		keys.Commands(appHomeDir),
	)

	flags.AddTxFlagsToCmd(rootCmd)

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

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
}
