package gorelayer

import (
	"fmt"
	wasmdtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	omniflixnfttypes "github.com/OmniFlix/onft/types"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
	rlycmd "github.com/cosmos/relayer/v2/cmd"
	"github.com/cosmos/relayer/v2/relayer"
	"github.com/cosmos/relayer/v2/relayer/chains/cosmos"
	"github.com/cosmos/relayer/v2/relayer/provider"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/spf13/viper"
	"github.com/strangelove-ventures/lens/client/codecs/ethermint"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
	"time"
)

func InitRly(log *zap.Logger) *Rly {
	rly := &Rly{
		Log:      log,
		Viper:    viper.New(),
		HomePath: filepath.Join(os.Getenv("HOME"), ".relayer"),
	}

	cfgPath := path.Join(rly.HomePath, "config", "config.yaml")
	rly.Viper.SetConfigFile(cfgPath)
	if err := rly.Viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// read the config file bytes
	file, err := os.ReadFile(rly.Viper.ConfigFileUsed())
	if err != nil {
		panic("Error reading file: " + err.Error())
	}

	// unmarshall them into the wrapper struct
	cfgWrapper := &rlycmd.ConfigInputWrapper{}
	err = yaml.Unmarshal(file, cfgWrapper)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshalling config: %v", err))
	}

	// verify that the channel filter rule is valid for every path in the config
	for _, p := range cfgWrapper.Paths {
		if err := p.ValidateChannelFilterRule(); err != nil {
			panic(fmt.Sprintf("error initializing the relayer config for path %s: %w", p.String(), err))
		}
	}

	// build the config struct
	chains := make(relayer.Chains)
	for chainName, pcfg := range cfgWrapper.ProviderConfigs {
		prov, err := pcfg.Value.(provider.ProviderConfig).NewProvider(
			rly.Log.With(zap.String("provider_type", pcfg.Type)),
			rly.HomePath, rly.Debug, chainName,
		)
		if err != nil {
			panic(fmt.Sprintf("failed to build ChainProviders: %w", err))
		}

		if err := prov.Init(); err != nil {
			panic(fmt.Sprintf("failed to initialize provider: %w", err))
		}

		chain := relayer.NewChain(rly.Log, prov, rly.Debug)
		chains[chainName] = chain
	}

	rly.Config = &rlycmd.Config{
		Global: cfgWrapper.Global,
		Chains: chains,
		Paths:  cfgWrapper.Paths,
	}

	// ensure config has []*relayer.Chain used for all chain operations
	if err := validateConfig(rly.Config); err != nil {
		panic(fmt.Sprintf("Error parsing chain config: %v", err))
	}

	for _, chain := range rly.Config.Chains {
		cosmosChain := chain.ChainProvider.(*cosmos.CosmosProvider)
		nfttransfertypes.RegisterInterfaces(cosmosChain.Codec.InterfaceRegistry)
		irisnfttypes.RegisterInterfaces(cosmosChain.Codec.InterfaceRegistry)
		irisnfttypes.RegisterLegacyAminoCodec(cosmosChain.Codec.Amino)

		omniflixnfttypes.RegisterInterfaces(cosmosChain.Codec.InterfaceRegistry)
		omniflixnfttypes.RegisterLegacyAminoCodec(cosmosChain.Codec.Amino)

		wasmdtypes.RegisterInterfaces(cosmosChain.Codec.InterfaceRegistry)
		wasmdtypes.RegisterLegacyAminoCodec(cosmosChain.Codec.Amino)

		ethermint.RegisterInterfaces(cosmosChain.Codec.InterfaceRegistry)
	}

	return rly
}

// validateConfig is used to validate the GlobalConfig values
func validateConfig(c *rlycmd.Config) error {
	_, err := time.ParseDuration(c.Global.Timeout)
	if err != nil {
		panic(fmt.Sprintf("did you remember to run 'rly config init' error: %v", err))
	}

	return nil
}
