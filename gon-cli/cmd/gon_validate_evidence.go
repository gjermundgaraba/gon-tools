package cmd

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"path"
	"strings"
)

func validateEvidenceFileInteractive() {
	relativePathToEvidence := askForString("Enter the relative path to the evidence file")

	if !strings.HasSuffix(".xlsx", relativePathToEvidence) {
		relativePathToEvidence = path.Join(relativePathToEvidence, "evidence.xlsx")
	}

	evidence, err := excelize.OpenFile(relativePathToEvidence)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Close the spreadsheet.
		if err := evidence.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	validateInfoSheet(evidence)

}

// ## Info
//
// - `sheet name` is unmodified (some people change the sheet name, but the validation tool has the sheet name hard-coded)
// - `row length` == `2` (some merge the row  2 and row 3, I get confused about this...)
// - input only in `row 2`
func validateInfoSheet(evidence *excelize.File) (validationErrors []string) {
	infoRows, err := evidence.GetRows("Info")
	if err != nil {
		return []string{"Info sheet not found"}
	}

	if len(infoRows) != 2 {
		validationErrors = append(validationErrors, "Info sheet should have 2 rows exactly")
	}

	return validationErrors
}
