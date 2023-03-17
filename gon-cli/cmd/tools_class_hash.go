package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	nfttransfertypes "github.com/bianjieai/nft-transfer/types"
)

func calculateClassHashInteractive() {
	trace := askForString("Enter the full class ibc trace (this is not validated in any way)", survey.WithValidator(survey.Required))

	fmt.Println("Class hash:")
	fmt.Println(calculateClassHash(trace))
}

func calculateClassHash(trace string) string {
	classTrace := nfttransfertypes.ParseClassTrace(trace)
	return fmt.Sprintf("ibc/%s", classTrace.Hash())
}
