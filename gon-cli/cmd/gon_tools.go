package cmd

import "github.com/spf13/cobra"

const (
	validateEvidenceFileOption  OptionString = "Validate Evidence File"
	validateEvidenceFileCommand              = "validate-evidence-file"
)

func gonToolsInteractive(cmd *cobra.Command, args []string) {
	toolsOptions := []OptionString{
		validateEvidenceFileOption,
	}

	var toolsChoice OptionString
	if len(args) > 1 && args[1] != "" {
		switch args[1] {
		case validateEvidenceFileCommand:
			toolsChoice = validateEvidenceFileOption
		default:
			panic("invalid command")
		}
	} else {
		toolsChoice = chooseOne("Which GoN tool do you want?", toolsOptions)
	}

	switch toolsChoice {
	case validateEvidenceFileOption:
		validateEvidenceFileInteractive()
	default:
		panic("invalid command")
	}
}
