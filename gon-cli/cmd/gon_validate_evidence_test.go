package cmd

import (
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
	"testing"
)

func createEmptySheetWithHeaders(evidence *excelize.File, sheet string, headers []string) {
	if _, err := evidence.NewSheet(sheet); err != nil {
		panic(err)
	}

	if err := evidence.SetSheetRow(sheet, "A1", &headers); err != nil {
		panic(err)
	}
}

func createEmptyInfoSheetWithHeaders(evidence *excelize.File) {
	createEmptySheetWithHeaders(evidence, "Info", []string{"TeamName", "IRISnetAddress", "StargazeAddress", "JunoAddress", "UptickAddress", "OmniFlixAddress", "DiscordHandle", "Community"})
}

func fillRowWithJunk(evidence *excelize.File, sheet, startingCell string, numberOfValues int) {
	row := make([]interface{}, numberOfValues)
	for i, _ := range row {
		row[i] = "test"
	}

	if err := evidence.SetSheetRow(sheet, startingCell, &row); err != nil {
		panic(err)
	}
}

func TestValidateInfoSheet(t *testing.T) {
	testTable := map[string]struct {
		setupTestEvidence func(evidence *excelize.File)
		expectedErrors    []string
	}{
		"Info sheet not found": {
			setupTestEvidence: func(evidence *excelize.File) {},
			expectedErrors:    []string{"Info sheet not found"},
		},
		"Info sheet should have 2 rows exactly (including the first header row)": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptyInfoSheetWithHeaders(evidence)
				fillRowWithJunk(evidence, "Info", "A2", 8)
				fillRowWithJunk(evidence, "Info", "A3", 8)
			},
			expectedErrors: []string{"Info sheet should have 2 rows exactly (including the first header row), but has 3 rows"},
		},
		"Info sheet headers should not be changed": {
			setupTestEvidence: func(evidence *excelize.File) {
				if _, err := evidence.NewSheet("Info"); err != nil {
					panic(err)
				}
				fillRowWithJunk(evidence, "Info", "A1", 8)
				fillRowWithJunk(evidence, "Info", "A2", 8)
			},
			expectedErrors: []string{
				"Sheet headers should not be changed, expected TeamName, got test",
				"Sheet headers should not be changed, expected IRISnetAddress, got test",
				"Sheet headers should not be changed, expected StargazeAddress, got test",
				"Sheet headers should not be changed, expected JunoAddress, got test",
				"Sheet headers should not be changed, expected UptickAddress, got test",
				"Sheet headers should not be changed, expected OmniFlixAddress, got test",
				"Sheet headers should not be changed, expected DiscordHandle, got test",
				"Sheet headers should not be changed, expected Community, got test"},
		},
		"Mandatory fields should not be empty": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptyInfoSheetWithHeaders(evidence)
				if err := evidence.SetCellStr("Info", "H2", "This one nobody cares about, the rest are mandatory"); err != nil {
					panic(err)
				}
			},
			expectedErrors: []string{
				"Column \"TeamName\", should not be empty",
				"Column \"IRISnetAddress\", should not be empty",
				"Column \"StargazeAddress\", should not be empty",
				"Column \"JunoAddress\", should not be empty",
				"Column \"UptickAddress\", should not be empty",
				"Column \"OmniFlixAddress\", should not be empty",
				"Column \"DiscordHandle\", should not be empty",
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			evidence := excelize.NewFile()
			test.setupTestEvidence(evidence)
			validationErrors := validateInfoSheet(evidence)
			require.Equal(t, test.expectedErrors, validationErrors)
		})
	}
}

func TestValidateA1(t *testing.T) {
	testTable := map[string]struct {
		setupTestEvidence func(evidence *excelize.File)
		expectedErrors    []string
	}{
		"A1 sheet not found": {
			setupTestEvidence: func(evidence *excelize.File) {},
			expectedErrors:    []string{"A1 sheet not found"},
		},
		"Info sheet should have 2 rows exactly (including the first header row)": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A1", []string{"TxHash", "ClassID"})
				fillRowWithJunk(evidence, "A1", "A2", 2)
				fillRowWithJunk(evidence, "A1", "A3", 2)
			},
			expectedErrors: []string{"A1 sheet should have 2 rows exactly (including the first header row), but has 3 rows"},
		},
		"A1 sheet headers should not be changed": {
			setupTestEvidence: func(evidence *excelize.File) {
				if _, err := evidence.NewSheet("A1"); err != nil {
					panic(err)
				}
				fillRowWithJunk(evidence, "A1", "A1", 2)
				fillRowWithJunk(evidence, "A1", "A2", 2)
			},
			expectedErrors: []string{
				"Sheet headers should not be changed, expected TxHash, got test",
				"Sheet headers should not be changed, expected ClassID, got test",
			},
		},
		"Mandatory fields should not be empty": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A1", []string{"TxHash", "ClassID"})
				if err := evidence.SetCellStr("A1", "A2", "mytxhash"); err != nil {
					panic(err)
				}
			},
			expectedErrors: []string{
				"All fields are mandatory, but row 2 has 1 column(s), expected 2",
			},
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			evidence := excelize.NewFile()
			test.setupTestEvidence(evidence)
			validationErrors := validateA1Sheet(evidence)
			require.Equal(t, test.expectedErrors, validationErrors)
		})
	}
}

func TestValidateA13(t *testing.T) {
	testTable := map[string]struct {
		setupTestEvidence func(evidence *excelize.File)
		expectedErrors    []string
	}{
		"A13 sheet not found": {
			setupTestEvidence: func(evidence *excelize.File) {},
			expectedErrors:    []string{"A13 sheet not found"},
		},
		"Info sheet should have 2 rows exactly (including the first header row)": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A13", []string{"TxHash", "ChainID"})
				fillRowWithJunk(evidence, "A13", "A2", 2)
				fillRowWithJunk(evidence, "A13", "A3", 2)
			},
			expectedErrors: []string{"A13 sheet should have 5 rows exactly (including the first header row), but has 3 rows"},
		},
		"A13 sheet headers should not be changed": {
			setupTestEvidence: func(evidence *excelize.File) {
				if _, err := evidence.NewSheet("A13"); err != nil {
					panic(err)
				}
				fillRowWithJunk(evidence, "A13", "A1", 2)

				if err := evidence.SetSheetRow("A13", "A2", &[]interface{}{"mytxhash", "gon-irishub-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A3", &[]interface{}{"mytxhash", "elgafar-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A4", &[]interface{}{"mytxhash", "uptick_7000-2"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A5", &[]interface{}{"mytxhash", "elgafar-1"}); err != nil {
					panic(err)
				}
			},
			expectedErrors: []string{
				"Sheet headers should not be changed, expected TxHash, got test",
				"Sheet headers should not be changed, expected ChainID, got test",
			},
		},
		"Mandatory fields should not be empty": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A13", []string{"TxHash", "ChainID"})
				if err := evidence.SetCellStr("A13", "A2", "mytxhash"); err != nil {
					panic(err)
				}
				if err := evidence.SetCellStr("A13", "A3", "mytxhash"); err != nil {
					panic(err)
				}
				if err := evidence.SetCellStr("A13", "A4", "mytxhash"); err != nil {
					panic(err)
				}
				if err := evidence.SetCellStr("A13", "A5", "mytxhash"); err != nil {
					panic(err)
				}

			},
			expectedErrors: []string{
				"All fields are mandatory, but row 2 has 1 column(s), expected 2",
				"All fields are mandatory, but row 3 has 1 column(s), expected 2",
				"All fields are mandatory, but row 4 has 1 column(s), expected 2",
				"All fields are mandatory, but row 5 has 1 column(s), expected 2",
			},
		},
		"Wrong chain id should fail": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A13", []string{"TxHash", "ChainID"})
				if err := evidence.SetSheetRow("A13", "A2", &[]interface{}{"7C0150BFDD1F6612A22E0B4A37B328FBF1EFD18A2C60F419DF04839C9A327BD3", "gon-irishub-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A3", &[]interface{}{"3CF3465E03060588D386CC7C02B64D6E620D336D2F438A20A581413E3A2AD9B3", "Elgafar-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A4", &[]interface{}{"1F935D8F74900FAB42E434DC047357426364AE85BE2453FDA54F65459FF20F42", "uptick_7000-2"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A5", &[]interface{}{"EECDF6ADA6691F8C458721E408139826C625C31C579B989FB9E0647F011FDFED", "elgafar-1"}); err != nil {
					panic(err)
				}
			},
			expectedErrors: []string{
				"ChainID for the row 3 should be \"elgafar-1\", but was \"Elgafar-1\"",
			},
		},
		"Valid sheet should be OK": {
			setupTestEvidence: func(evidence *excelize.File) {
				createEmptySheetWithHeaders(evidence, "A13", []string{"TxHash", "ChainID"})
				if err := evidence.SetSheetRow("A13", "A2", &[]interface{}{"7C0150BFDD1F6612A22E0B4A37B328FBF1EFD18A2C60F419DF04839C9A327BD3", "gon-irishub-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A3", &[]interface{}{"3CF3465E03060588D386CC7C02B64D6E620D336D2F438A20A581413E3A2AD9B3", "elgafar-1"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A4", &[]interface{}{"1F935D8F74900FAB42E434DC047357426364AE85BE2453FDA54F65459FF20F42", "uptick_7000-2"}); err != nil {
					panic(err)
				}
				if err := evidence.SetSheetRow("A13", "A5", &[]interface{}{"EECDF6ADA6691F8C458721E408139826C625C31C579B989FB9E0647F011FDFED", "elgafar-1"}); err != nil {
					panic(err)
				}
			},
			expectedErrors: nil,
		},
	}

	for name, test := range testTable {
		t.Run(name, func(t *testing.T) {
			evidence := excelize.NewFile()
			test.setupTestEvidence(evidence)
			validationErrors := validateA13Sheet(evidence)
			require.Equal(t, test.expectedErrors, validationErrors)
		})
	}
}
