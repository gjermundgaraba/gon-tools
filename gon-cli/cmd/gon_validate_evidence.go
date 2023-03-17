package cmd

import (
	"fmt"
	"github.com/gjermundgaraba/gon/chains"
	"github.com/xuri/excelize/v2"
	"path"
	"strconv"
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

	errors := validateEvidenceFile(evidence)
	errorsFound := false
	for sheet, sheetErrors := range errors {
		if len(sheetErrors) > 0 {
			errorsFound = true
			fmt.Println()
			fmt.Println("Errors in sheet", sheet)
			for _, sheetError := range sheetErrors {
				fmt.Println(sheetError)
			}
		}
	}

	if !errorsFound {
		fmt.Println("Congratulations! Your evidence file appears to be valid!")
		fmt.Println("Keep in mind, not all fields are fully validated here yet, so please double check your evidence file.")
	}
}

func validateEvidenceFile(evidence *excelize.File) (validationErrors map[string][]string) {
	validationErrors = make(map[string][]string)

	validationErrors["Info"] = validateInfoSheet(evidence)
	validationErrors["A1"] = validateA1Sheet(evidence)
	validationErrors["A2"] = validateA2Sheet(evidence)
	validationErrors["A3"] = validateA3Sheet(evidence)
	validationErrors["A4"] = validateA4Sheet(evidence)
	validationErrors["A5"] = validateA5Sheet(evidence)
	validationErrors["A6"] = validateA6Sheet(evidence)
	validationErrors["A7"] = validateA7Sheet(evidence)
	validationErrors["A8"] = validateA8Sheet(evidence)
	validationErrors["A9"] = validateA9Sheet(evidence)
	validationErrors["A10"] = validateA10Sheet(evidence)
	validationErrors["A11"] = validateA11Sheet(evidence)
	validationErrors["A12"] = validateA12Sheet(evidence)
	validationErrors["A13"] = validateA13Sheet(evidence)
	validationErrors["A14"] = validateA14Sheet(evidence)
	validationErrors["A15"] = validateA15Sheet(evidence)
	validationErrors["A16"] = validateA16Sheet(evidence)
	validationErrors["A17"] = validateA17Sheet(evidence)
	validationErrors["A18"] = validateA18Sheet(evidence)
	validationErrors["A19"] = validateA19Sheet(evidence)
	validationErrors["A20"] = validateA20Sheet(evidence)

	return validationErrors
}

// ## Info
//
// - `sheet name` is unmodified (some people change the sheet name, but the validation tool has the sheet name hard-coded)
// - `row length` == `2` (some merge the row  2 and row 3, I get confused about this...)
// - input only in `row 2`
func validateInfoSheet(evidence *excelize.File) (validationErrors []string) {
	validationErrors = validateSheet(evidence, "Info", 2, []string{"TeamName", "IRISnetAddress", "StargazeAddress", "JunoAddress", "UptickAddress", "OmniFlixAddress", "DiscordHandle", "Community"}, false)
	if len(validationErrors) > 0 {
		return validationErrors
	}
	infoRows, _ := evidence.GetRows("Info")

	for i, _ := range make([]int, 7) {
		if infoRows[1][i] == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("Column %q, should not be empty", infoRows[0][i]))
		}
	}

	return validationErrors
}

// ## A1
//
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA1Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A1", 2, []string{"TxHash", "ClassID"}, true)
}

// ## A2
//
// - `sheet name` is unmodified
// - `row length` >= `3`
// - input from `row 2` (but we check only the first 2 NFTs)
func validateA2Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A2", 3, []string{"TxHash", "ClassID", "NFTID"}, true)
}

// ## A3
//
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA3Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A3", 2, []string{"TxHash", "ClassID", "NFTID", "ChainID"}, true)
}

// ## A4
//
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA4Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A4", 2, []string{"TxHash", "ClassID", "NFTID", "ChainID"}, true)
}

// ## A5
//
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA5Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A5", 2, []string{"TxHash", "ClassID", "NFTID", "ChainID"}, true)
}

// ## A6
//
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA6Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A6", 2, []string{"TxHash", "ClassID", "NFTID", "ChainID"}, true)
}

// ## A7
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA7Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A7", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A8
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA8Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A8", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A9
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA9Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A9", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A10
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA10Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A10", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A11
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA11Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A11", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A12
//
// for each of these tasks:
// - `sheet name` is unmodified
// - `row length` == `2`
// - input only in `row 2`
func validateA12Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheet(evidence, "A12", 2, []string{"ClassID", "NFTID"}, true)
}

// ## A13
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA13Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A13", []string{"TxHash", "ChainID"}, true, "i --(1)--> s --(1)--> u --(1)--> s --(2)--> i")
}

// ## A14
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA14Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A14", []string{"TxHash", "ChainID"}, true, "i --(1)--> u --(1)--> o --(1)--> u --(2)--> i")
}

// ## A15
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA15Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A15", []string{"TxHash", "ChainID"}, true, "i --(1)--> j --(1)--> u --(1)--> j --(2)--> i")
}

// ## A16
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA16Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A16", []string{"TxHash", "ChainID"}, true, "i --(1)--> j --(1)--> s --(1)--> j --(2)--> i")
}

// ## A17
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA17Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A17", []string{"TxHash", "ChainID"}, true, "i --(1)--> s --(1)--> j --(1)--> s --(1)--> i")
}

// ## A18
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA18Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A18", []string{"TxHash", "ChainID"}, true, "i --(1)--> o --(1)--> u --(1)--> o --(1)--> i")
}

// ## A19
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA19Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A19", []string{"TxHash", "ChainID"}, true, "i --(1)--> s --(1)--> j --(1)--> u --(1)--> j --(1)--> s --(1)--> i")
}

// ## A20
//
// - `sheet name` is unmodified
// - `row length` == `max-hop + 1`  (for example: in `i -> s -> j -> i`, `max-hop` is `3` and `row length` must be `4`)
// - input from `row 2`
func validateA20Sheet(evidence *excelize.File) (validationErrors []string) {
	return validateSheetWithFlow13to20(evidence, "A20", []string{"TxHash", "ChainID"}, true, "i --(1)--> u --(1)--> s --(1)--> o --(1)--> s --(1)--> u --(1)--> i")
}

func validateSheetWithFlow13to20(evidence *excelize.File, sheetName string, expectedHeaders []string, allFieldsMandatory bool, flow string) (validationErrors []string) {
	connections := parseFlow(flow)
	expectedNumberOfRows := len(connections) + 1
	validationErrors = validateSheet(evidence, sheetName, expectedNumberOfRows, expectedHeaders, allFieldsMandatory)
	if len(validationErrors) > 0 {
		return validationErrors
	}

	// check if the flow is correct
	for i, connection := range connections {
		value, err := evidence.GetCellValue(sheetName, fmt.Sprintf("B%d", i+2))
		if err != nil {
			panic(err)
		}

		if value != string(connection.ChannelA.ChainID) {
			validationErrors = append(validationErrors, fmt.Sprintf("ChainID for the row %d should be %q, but was %q", i+2, connection.ChannelA.ChainID, value))
		}
	}

	return validationErrors
}

var chainMap = map[string]chains.Chain{
	"i": chains.IRISChain,
	"s": chains.StargazeChain,
	"j": chains.JunoChain,
	"o": chains.OmniFlixChain,
	"u": chains.UptickChain,
}

// example: i --(1)--> s --(1)--> j --(1)--> i
func parseFlow(flow string) []chains.NFTConnection {
	var connections []chains.NFTConnection
	flowPieces := strings.Split(flow, ")-->")
	for i, fp := range flowPieces {
		if i != len(flowPieces)-1 {
			chain := chainMap[string(strings.TrimSpace(fp)[0])]
			nextChain := chainMap[string(strings.TrimSpace(flowPieces[i+1])[0])]
			connectionNumber, err := strconv.ParseInt(string(fp[len(fp)-1]), 10, 64)
			if err != nil {
				panic(err)
			}
			connections = append(connections, chain.GetConnectionsTo(nextChain)[connectionNumber-1])
		}
	}

	return connections
}

func validateSheet(evidence *excelize.File, sheetName string, expectedNumberOfRows int, expectedHeaders []string, allFieldsMandatory bool) (validationErrors []string) {
	rows, err := evidence.GetRows(sheetName)
	if err != nil {
		return []string{fmt.Sprintf("%s sheet not found", sheetName)}
	}

	if len(rows) != expectedNumberOfRows {
		return append(validationErrors, fmt.Sprintf("%s sheet should have %d rows exactly (including the first header row), but has %d rows", sheetName, expectedNumberOfRows, len(rows)))
	}

	for i, header := range expectedHeaders {
		actualHeader := rows[0][i]
		if actualHeader != header {
			validationErrors = append(validationErrors, fmt.Sprintf("Sheet headers should not be changed, expected %s, got %s", header, actualHeader))
		}
	}

	if allFieldsMandatory {
		// Skip the first row, which is the header
		for i := 1; i < expectedNumberOfRows; i++ {
			if len(rows[i]) != len(expectedHeaders) {
				validationErrors = append(validationErrors, fmt.Sprintf("All fields are mandatory, but row %d has %d column(s), expected %d", i+1, len(rows[i]), len(expectedHeaders)))
				continue
			}
			for j, cell := range rows[i] {
				if cell == "" {
					validationErrors = append(validationErrors, fmt.Sprintf("All fields are mandatory, but %s is empty", expectedHeaders[j]))
				}
			}
		}
	}

	return validationErrors
}
