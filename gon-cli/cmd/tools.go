package cmd

import "github.com/spf13/cobra"

const (
	findIBCTransactionsOption  OptionString = "Find IBC Transactions"
	findIBCTransactionsCommand              = "find-ibc-txs"

	listConnectionsOption  OptionString = "List Connections"
	listConnectionsCommand              = "list-connections"

	manageKeysOption  OptionString = "Manage Keys"
	manageKeysCommand              = "manage-keys"

	queryTransactionOption  OptionString = "Query Transaction"
	queryTransactionCommand              = "query-tx"

	calculateClassHashOption  OptionString = "Calculate Class Hash"
	calculateClassHashCommand              = "class-hash"

	generateTraceOption  OptionString = "Generate Trace manually"
	generateTraceCommand              = "generate-trace"

	generateTraceFromFlowOption  OptionString = "Generate Trace manually from flow"
	generateTraceFromFlowCommand              = "generate-trace-from-flow"
)

func toolsInteractive(cmd *cobra.Command, args []string) {
	toolsOptions := []OptionString{
		findIBCTransactionsOption,
		listConnectionsOption,
		manageKeysOption,
		queryTransactionOption,
		calculateClassHashOption,
		generateTraceOption,
		generateTraceFromFlowOption,
	}

	var toolsChoice OptionString
	if len(args) > 1 && args[1] != "" {
		switch args[1] {
		case findIBCTransactionsCommand:
			toolsChoice = findIBCTransactionsOption
		case listConnectionsCommand:
			toolsChoice = listConnectionsOption
		case manageKeysCommand:
			toolsChoice = manageKeysOption
		case queryTransactionCommand:
			toolsChoice = queryTransactionOption
		case calculateClassHashCommand:
			toolsChoice = calculateClassHashOption
		case generateTraceCommand:
			toolsChoice = generateTraceOption
		case generateTraceFromFlowCommand:
			toolsChoice = generateTraceFromFlowOption
		default:
			panic("invalid command")
		}
	} else {
		toolsChoice = chooseOne("Which tool do you want?", toolsOptions)
	}

	switch toolsChoice {
	case findIBCTransactionsOption:
		findIBCTransactionsInteractive(cmd)
	case listConnectionsOption:
		listConnections()
	case manageKeysOption:
		manageKeys(cmd)
	case queryTransactionOption:
		queryTransaction(cmd)
	case calculateClassHashOption:
		calculateClassHashInteractive()
	case generateTraceOption:
		generateTraceInteractive()
	case generateTraceFromFlowOption:
		generateTraceSimplyInteractive()
	default:
		panic("invalid command")
	}
}
