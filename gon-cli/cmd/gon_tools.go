package cmd

import "github.com/spf13/cobra"

const (
	validateEvidenceFileOption  OptionString = "Validate Evidence File"
	validateEvidenceFileCommand              = "validate-evidence-file"

	raceOption  OptionString = "Race"
	raceCommand              = "race"

	quizOption  OptionString = "Quiz"
	quizCommand              = "quiz"
)

func gonToolsInteractive(cmd *cobra.Command, args []string, appHomeDir string) {
	toolsOptions := []OptionString{
		validateEvidenceFileOption,
		raceOption,
		quizOption,
	}

	var toolsChoice OptionString
	if len(args) > 1 && args[1] != "" {
		switch args[1] {
		case validateEvidenceFileCommand:
			toolsChoice = validateEvidenceFileOption
		case raceCommand:
			toolsChoice = raceOption
		case quizCommand:
			toolsChoice = quizOption
		default:
			panic("invalid command")
		}
	} else {
		toolsChoice = chooseOne("Which GoN tool do you want?", toolsOptions)
	}

	switch toolsChoice {
	case validateEvidenceFileOption:
		validateEvidenceFileInteractive()
	case raceOption:
		raceInteractive(cmd, appHomeDir)
	case quizOption:
		gonQuizInteractive()
	default:
		panic("invalid command")
	}
}
