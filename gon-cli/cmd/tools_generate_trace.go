package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gjermundgaraba/gon/chains"
	"strings"
)

func generateTraceInteractive() {
	className := askForString("Enter the original class name", survey.WithValidator(survey.Required))
	currentChain := chooseChain("Choose the initial source chain")

	trace := className
	for {
		destinationChain := chooseChain("Choose the next chain", currentChain)
		connection := chooseConnection(currentChain, destinationChain, "Choose the connection to use")
		trace, _ = calculateClassTrace(trace, connection)

		currentChain = destinationChain

		if anotherHop := askForConfirmation("Add another hop?", true); anotherHop {
			continue
		}
		break
	}

	fmt.Println("Full IBC class trace:")
	fmt.Println(trace)
	if len(strings.Split(trace, "/")) > 2 && currentChain.NFTImplementation() == chains.CosmosSDK {
		fmt.Println()
		fmt.Println("Class hash:")
		fmt.Println(calculateClassHash(trace))
	}
}
