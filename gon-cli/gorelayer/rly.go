package gorelayer

import (
	"fmt"
	rlycmd "github.com/cosmos/relayer/v2/cmd"
	"github.com/cosmos/relayer/v2/relayer"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var pathMap = map[string]map[string]string{
	"gon-flixnet-1": {
		"channel-24": "gon-irishub-1_channel-0-gon-flixnet-1_channel-24",
		"channel-25": "gon-irishub-1_channel-1-gon-flixnet-1_channel-25",
		"channel-41": "uptick_7000-2_channel-5-gon-flixnet-1_channel-41",
		"channel-42": "uptick_7000-2_channel-9-gon-flixnet-1_channel-42",
		"channel-44": "gon-flixnet-1_channel-44-elgafar-1_channel-209",
		"channel-45": "gon-flixnet-1_channel-45-elgafar-1_channel-210",
		"channel-46": "gon-flixnet-1_channel-46-uni-6_channel-91",
		"channel-47": "gon-flixnet-1_channel-47-uni-6_channel-92",
	},
	"gon-irishub-1": {
		"channel-0":  "gon-irishub-1_channel-0-gon-flixnet-1_channel-24",
		"channel-1":  "gon-irishub-1_channel-1-gon-flixnet-1_channel-25",
		"channel-17": "gon-irishub-1_channel-17-uptick_7000-2_channel-3",
		"channel-19": "gon-irishub-1_channel-19-uptick_7000-2_channel-4",
		"channel-22": "gon-irishub-1_channel-22-elgafar-1_channel-207",
		"channel-23": "gon-irishub-1_channel-23-elgafar-1_channel-208",
		"channel-24": "gon-irishub-1_channel-24-uni-6_channel-89",
		"channel-25": "gon-irishub-1_channel-25-uni-6_channel-90",
	},
	"uni-6": {
		"channel-86":  "uptick_7000-2_channel-7-uni-6_channel-86",
		"channel-88":  "uptick_7000-2_channel-13-uni-6_channel-88",
		"channel-89":  "gon-irishub-1_channel-24-uni-6_channel-89",
		"channel-90":  "gon-irishub-1_channel-25-uni-6_channel-90",
		"channel-91":  "gon-flixnet-1_channel-46-uni-6_channel-91",
		"channel-92":  "gon-flixnet-1_channel-47-uni-6_channel-92",
		"channel-93":  "uni-6_channel-93-elgafar-1_channel-211",
		"channel-94":  "uni-6_channel-94-elgafar-1_channel-213",
		"channel-120": "uni-6_channel-120-elgafar-1_channel-230",
		"channel-122": "uni-6_channel-122-elgafar-1_channel-234",
	},
	"uptick_7000-2": {
		"channel-3":  "gon-irishub-1_channel-17-uptick_7000-2_channel-3",
		"channel-4":  "gon-irishub-1_channel-19-uptick_7000-2_channel-4",
		"channel-5":  "uptick_7000-2_channel-5-gon-flixnet-1_channel-41",
		"channel-6":  "uptick_7000-2_channel-6-elgafar-1_channel-203",
		"channel-7":  "uptick_7000-2_channel-7-uni-6_channel-86",
		"channel-9":  "uptick_7000-2_channel-9-gon-flixnet-1_channel-42",
		"channel-12": "uptick_7000-2_channel-12-elgafar-1_channel-206",
		"channel-13": "uptick_7000-2_channel-13-uni-6_channel-88",
	},
	"elgafar-1": {
		"channel-203": "uptick_7000-2_channel-6-elgafar-1_channel-203",
		"channel-206": "uptick_7000-2_channel-12-elgafar-1_channel-206",
		"channel-207": "gon-irishub-1_channel-22-elgafar-1_channel-207",
		"channel-208": "gon-irishub-1_channel-23-elgafar-1_channel-208",
		"channel-209": "gon-flixnet-1_channel-44-elgafar-1_channel-209",
		"channel-210": "gon-flixnet-1_channel-45-elgafar-1_channel-210",
		"channel-211": "uni-6_channel-93-elgafar-1_channel-211",
		"channel-213": "uni-6_channel-94-elgafar-1_channel-213",
		"channel-230": "uni-6_channel-120-elgafar-1_channel-230",
		"channel-234": "uni-6_channel-122-elgafar-1_channel-234",
	},
}

type Rly struct {
	// Log is the root logger of the application.
	// Consumers are expected to store and use local copies of the logger
	// after modifying with the .With method.
	Log *zap.Logger

	Viper *viper.Viper

	HomePath string
	Debug    bool
	Config   *rlycmd.Config
}

func (rly *Rly) GetRelayerChain(gonChain chains.Chain) *relayer.Chain {
	chain, err := rly.Config.Chains.Get(string(gonChain.ChainID()))
	if err != nil {
		panic(err)
	}

	return chain
}

func (rly *Rly) GetPathString(connection chains.NFTConnection) string {
	path, ok := pathMap[string(connection.ChannelA.ChainID)][connection.ChannelA.Channel]
	if !ok {
		panic(fmt.Sprintf("Path not found for chainID %s and channel %s", connection.ChannelA.ChainID, connection.ChannelA.Channel))
	}

	return path
}

func (rly *Rly) GetPath() {

}
