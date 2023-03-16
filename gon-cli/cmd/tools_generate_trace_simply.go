package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func generateTraceSimplyInteractive() {
	className := askForString("Enter the original class name", survey.WithValidator(survey.Required))
	simplyPath := askForString("Enter path (ie. i1j2u)", survey.WithValidator(survey.Required))

	simplyPath = pathReducer2(simplyPath)

	PORT_CHANNEL_MATRIX := map[string]string{
		"01": "seg1",
		"02": "seg2",
		"03": "seg3",
		// add more entries as needed
	}

	fullPath := ""
	for i := 0; i < len(simplyPath); i++ {
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

// func pathReducer(path string) string {
// 	newPath := ""
// 	pal := ""
// 	piece := ""
// 	letter := ""
// 	channel := ""
// 	fLetter := ""

// 	for i := 0; i < len(path); i += 3 {
// 		piece = path[i:2]
// 		letter = piece[0:1]
// 		channel = piece[1:2]
// 		if channel == "2" {
// 			letter = strings.ToUpper(letter)
// 		}
// 		newPath = newPath + letter
// 	}

// 	pal = findPalindrome(newPath)
// 	for pal != "" {
// 		fLetter = pal[0:1]
// 		newPath = strings.Replace(newPath, pal, fLetter, -1)
// 		pal = findPalindrome(newPath)
// 	}

// 	return newPath
// }

// func findPalindrome(path string) string {
// 	length := len(path)
// 	reverseSubstring := ""
// 	substring := ""
// 	minPalindrome := ""
// 	for i := 0; i < length; i++ {
// 		for j := 2; j <= length-i; j++ {
// 			substring = path[i:j]
// 			reverseSubstring = reverseString(substring)
// 			if substring == reverseSubstring {
// 				minPalindrome = substring
// 			}
// 		}
// 	}
// 	return minPalindrome
// }

// func reverseString(str string) string {
// 	byte_str := []rune(str)
// 	for i, j := 0, len(byte_str)-1; i < j; i, j = i+1, j-1 {
// 		byte_str[i], byte_str[j] = byte_str[j], byte_str[i]
// 	}
// 	return string(byte_str)
// }

func pathReducer2(input string) string {
	newPath := ""
	for i := 0; i < len(input); i += 2 {
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
