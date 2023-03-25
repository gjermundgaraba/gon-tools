package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/gjermundgaraba/gon/gorelayer"
	irisnfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

type airdropData struct {
	Type        string `json:"type"`
	Flow        string `json:"flow"`
	LastOwner   string `json:"last_owner"`
	StartHeight string `json:"start_height"`
}

type raceFlow struct {
	FlowName string
	FlowRaw  string
	Path     []chains.NFTConnection
}

var raceFlows = map[string]raceFlow{
	"f01": {
		FlowName: "f1",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f1": {
		FlowName: "f1",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f02": {
		FlowName: "f2",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f2": {
		FlowName: "f2",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f03": {
		FlowName: "f3",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f3": {
		FlowName: "f3",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f04": {
		FlowName: "f4",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f4": {
		FlowName: "f4",
		FlowRaw:  "i --(1)--> s --(1)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f05": {
		FlowName: "f5",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f5": {
		FlowName: "f5",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f06": {
		FlowName: "f6",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f6": {
		FlowName: "f6",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f07": {
		FlowName: "f7",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f7": {
		FlowName: "f7",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f08": {
		FlowName: "f8",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f8": {
		FlowName: "f8",
		FlowRaw:  "i --(1)--> s --(1)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f09": {
		FlowName: "f9",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f9": {
		FlowName: "f9",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f010": {
		FlowName: "f10",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f10": {
		FlowName: "f10",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f011": {
		FlowName: "f11",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f11": {
		FlowName: "f11",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f012": {
		FlowName: "f12",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f12": {
		FlowName: "f12",
		FlowRaw:  "i --(1)--> s --(2)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f013": {
		FlowName: "f13",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f13": {
		FlowName: "f13",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f014": {
		FlowName: "f14",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f14": {
		FlowName: "f14",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f015": {
		FlowName: "f15",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f15": {
		FlowName: "f15",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f016": {
		FlowName: "f16",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f16": {
		FlowName: "f16",
		FlowRaw:  "i --(1)--> s --(2)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[0],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f017": {
		FlowName: "f17",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f17": {
		FlowName: "f17",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f018": {
		FlowName: "f18",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f18": {
		FlowName: "f18",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f019": {
		FlowName: "f19",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f19": {
		FlowName: "f19",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f020": {
		FlowName: "f20",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f20": {
		FlowName: "f20",
		FlowRaw:  "i --(2)--> s --(1)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f021": {
		FlowName: "f21",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f21": {
		FlowName: "f21",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f022": {
		FlowName: "f22",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f22": {
		FlowName: "f22",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f023": {
		FlowName: "f23",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f23": {
		FlowName: "f23",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f024": {
		FlowName: "f24",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f24": {
		FlowName: "f24",
		FlowRaw:  "i --(2)--> s --(1)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[0],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f025": {
		FlowName: "f25",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f25": {
		FlowName: "f25",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f026": {
		FlowName: "f26",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f26": {
		FlowName: "f26",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f027": {
		FlowName: "f27",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f27": {
		FlowName: "f27",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f028": {
		FlowName: "f28",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f28": {
		FlowName: "f28",
		FlowRaw:  "i --(2)--> s --(2)--> j --(1)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[0],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f029": {
		FlowName: "f29",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f29": {
		FlowName: "f29",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(1)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f030": {
		FlowName: "f30",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f30": {
		FlowName: "f30",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(1)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[0],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f031": {
		FlowName: "f31",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f31": {
		FlowName: "f31",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(2)--> o --(1)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[0],
		},
	},
	"f032": {
		FlowName: "f32",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
	"f32": {
		FlowName: "f32",
		FlowRaw:  "i --(2)--> s --(2)--> j --(2)--> u --(2)--> o --(2)--> i",
		Path: []chains.NFTConnection{
			chains.IRISChain.GetConnectionsTo(chains.StargazeChain)[1],
			chains.StargazeChain.GetConnectionsTo(chains.JunoChain)[1],
			chains.JunoChain.GetConnectionsTo(chains.UptickChain)[1],
			chains.UptickChain.GetConnectionsTo(chains.OmniFlixChain)[1],
			chains.OmniFlixChain.GetConnectionsTo(chains.IRISChain)[1],
		},
	},
}

func raceInteractive(cmd *cobra.Command, appHomeDir string) {
	iris := chains.IRISChain
	key := chooseOrCreateKey(cmd, iris)
	ethKey := getEthermintKeyName(key)
	irisQueryClientCtx := getQueryClientContext(cmd, iris)
	irisAddress := getAddressForChain(cmd, iris, key)

	verbose, _ := cmd.Flags().GetBool(flagVerbose)
	kr := getKeyring(cmd)
	rly := gorelayer.InitRly(appHomeDir, key, ethKey, kr, verbose)

	selectedClass := getUsersNfts(cmd.Context(), getQueryClientContext(cmd, iris), iris, irisAddress)
	if len(selectedClass.NFTs) == 0 {
		fmt.Println("No NFT classes found")
		return
	}

	selectedNFT := chooseOne("Select NFT", selectedClass.NFTs)

	nftQueryClient := irisnfttypes.NewQueryClient(irisQueryClientCtx)
	nft, err := nftQueryClient.NFT(cmd.Context(), &irisnfttypes.QueryNFTRequest{
		DenomId: selectedClass.ClassID,
		TokenId: selectedNFT.ID,
	})
	if err != nil {
		panic(err)
	}

	var nftData airdropData
	fmt.Println("NFT data:", nft.NFT.Data)
	err = json.Unmarshal([]byte(nft.NFT.Data), &nftData)

	height, _ := getCurrentChainStatus(cmd.Context(), irisQueryClientCtx)
	startHeight, err := strconv.ParseUint(nftData.StartHeight, 10, 64)
	if err != nil {
		panic(err)
	}
	if height < startHeight {
		fmt.Printf("NFT start height has not happened yet. You need to wait! Current height: %d, NFT start height: %s\n", height, nftData.StartHeight)
		return
	}

	flow := raceFlows[nftData.Flow]
	fmt.Println("Selected flow:", flow.FlowName)
	fmt.Println("Raw flow:", flow.FlowRaw)
	finalClass := runFlow(cmd, rly, key, flow, selectedClass, selectedNFT)

	setAddressPrefixes(iris.Bech32Prefix())
	irisTxClientCtx := getClientTxContext(cmd, iris).
		WithSkipConfirmation(true)

	msgTransferNFT := irisnfttypes.NewMsgTransferNFT(selectedNFT.ID, finalClass.ClassID, "[do-not-modify]", "[do-not-modify]", "[do-not-modify]", "[do-not-modify]", irisAddress, nftData.LastOwner)
	fmt.Println("Message to be sent:")
	_ = irisTxClientCtx.PrintProto(msgTransferNFT)
	txResponse, err := sendTX(irisTxClientCtx, cmd.Flags(), msgTransferNFT)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tx hash of final transfer message (before success):", txResponse.TxHash)
	waitForTX(cmd, iris, txResponse.TxHash, "Final transfer to last_owner", "Final transfer to last_owner")
}

func runFlow(cmd *cobra.Command, rly *gorelayer.Rly, key string, flow raceFlow, selectedClass chains.NFTClass, selectedNFT chains.NFT) chains.NFTClass {
	currentClass := selectedClass
	for _, conn := range flow.Path {
		fmt.Println()
		fmt.Printf("Sending NFT from %s (ch: %s) to %s (ch: %s)\n", conn.ChannelA.ChainID, conn.ChannelA.Channel, conn.ChannelB.ChainID, conn.ChannelB.Channel)
		fmt.Printf("Current NFT class: %s\n", currentClass.ClassID)
		fmt.Printf("Current NFT trace: %s\n", currentClass.FullPathClassID)
		sourceChain := chains.GetChainFromChainID(conn.ChannelA.ChainID)
		destinationChain := chains.GetChainFromChainID(conn.ChannelB.ChainID)
		currentClass = transferNFT(cmd, rly, key, sourceChain, destinationChain, conn, currentClass, selectedNFT)
		time.Sleep(500 * time.Millisecond)
	}

	return currentClass
}

func transferNFT(cmd *cobra.Command, rly *gorelayer.Rly, keyName string, sourceChain, destinationChain chains.Chain, connection chains.NFTConnection, class chains.NFTClass, nft chains.NFT) chains.NFTClass {
	if err := cmd.Flags().Set(flags.FlagFrom, getCorrectedKeyName(keyName, sourceChain)); err != nil {
		panic(err)
	}
	setAddressPrefixes(sourceChain.Bech32Prefix())
	fromAddress := getAddressForChain(cmd, sourceChain, keyName)
	toAddress := getAddressForChain(cmd, destinationChain, keyName)

	fmt.Println("From address:", fromAddress)
	fmt.Println("To address:", toAddress)

	clientCtx := getClientTxContext(cmd, sourceChain).
		WithSkipConfirmation(true)
	targetChainHeight, targetChainTimestamp := getCurrentChainStatus(cmd.Context(), getQueryClientContext(cmd, destinationChain))
	timeoutHeight, timeoutTimestamp := sourceChain.GetIBCTimeouts(clientCtx, connection.ChannelA.Port, connection.ChannelA.Channel, targetChainHeight, targetChainTimestamp, false)
	msg := sourceChain.CreateTransferNFTMsg(connection.ChannelA, class, nft, fromAddress, toAddress, timeoutHeight, timeoutTimestamp)

	fmt.Println("Message to be sent:")
	_ = clientCtx.PrintProto(msg)

	txResponse, err := sendTX(clientCtx, cmd.Flags(), msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tx hash (before success):", txResponse.TxHash)

	expectedDestinationTrace, _ := calculateClassTrace(class.FullPathClassID, connection)
	fmt.Println("Expected destination trace:", expectedDestinationTrace)

	waitAndPrintIBCTrail(cmd, sourceChain, destinationChain, txResponse.TxHash, rly, false)

	return queryNftClassFromTrace(cmd, expectedDestinationTrace, destinationChain)
}
