package cmd

import (
	"fmt"
	"strings"

	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

var cleanPath = regexp.MustCompile(`[^ijous12]+`)

func generateTraceSimplyInteractive() {
	className := askForString("Enter the original class name", survey.WithValidator(survey.Required))
	simplyPath := askForString("Enter path (ie. i1j2u or i --(1)--> j --(2)--> u)", survey.WithValidator(survey.Required))
	simplyPath = cleanPath.ReplaceAllString(simplyPath, "")
	simplyPath = pathReducer(simplyPath)
	PORT_CHANNEL_MATRIX := getMatrix()
	fullPath := ""
	for i := 0; i < len(simplyPath); i++ {
		if len(simplyPath) < (i + 2) {
			break
		}
		ch := simplyPath[i : i+2]
		seg, ok := PORT_CHANNEL_MATRIX[ch]
		if ok {
			fullPath = seg + "/" + fullPath
		}
	}

	trace := fullPath + className

	fmt.Println("Full IBC class trace:")
	fmt.Println(trace)
	/*if len(strings.Split(trace, "/")) > 2 && currentChain.NFTImplementation() == chains.CosmosSDK {
		fmt.Println()
		fmt.Println("Class hash:")
		fmt.Println(calculateClassHash(trace))
	}*/
}

func pathReducer(input string) string {
	newPath := ""
	for i := 0; i < len(input); i += 2 {
		if len(input) < (i + 2) {
			newPath += input[i : i+1]
			break
		}
		piece := input[i : i+2]
		letter := string(piece[0])
		channel := string(piece[1])
		if channel == "2" {
			letter = strings.ToUpper(letter)
		}

		newPath += letter
	}

	reverse := func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	findPalindrome := func(path string) string {
		length := len(path)
		minPalindrome := ""

		for i := 0; i < length; i++ {
			for j := 2; j <= length-i; j++ {
				substring := path[i : i+j]
				reverseSubstring := reverse(substring)

				if substring == reverseSubstring {
					if minPalindrome == "" || len(substring) < len(minPalindrome) {
						minPalindrome = substring
					}
				}
			}
		}

		return minPalindrome
	}

	pal := findPalindrome(newPath)

	for pal != "" {
		fLetter := string(pal[0])
		newPath = strings.ReplaceAll(newPath, pal, fLetter)
		pal = findPalindrome(newPath)
	}

	return newPath
}
