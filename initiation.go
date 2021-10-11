package main

import (
	"fmt"
)

// getMode receives the desired gossip method which the user wants
func getMode() string {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Welcome! Which method of gossip would you like to implement: " +
		"Push(1), Pull(2), or Push/Pull Original(3) or Push/Pull Switch(4)? Please enter the number code as indicated.")
	fmt.Println("If you would like to quit, please enter 'q'. ")
	_, err := fmt.Scanf("%s", &mode)
	if err != nil || isValidMode(mode) {
		fmt.Println("Invalid Code: Please select a valid gossip protocol next time!")
		return "QUIT"
	} else if mode == "q" || mode == "Q" {
		fmt.Println("You have quit the program. Thank you and goodbye!")
		return "QUIT"
	}
	return mode
}

// isValidMode returns a boolean value signifying the user's choice of gossip mode can be accepted by the program
func isValidMode(mode string) bool {
	if mode != "1" && mode != "2" && mode != "3" && mode != "4" && mode != "q" && mode != "Q" {
		return true
	}
	return false
}

// getSettings returns the maximum amount of nodes and the user's choice of verbose output
func getSettings() (int, string) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("For up to how many nodes would you like to test the convergence speed?")
	_, err1 := fmt.Scanf("%d", &desiredNodes)
	if err1 != nil {
		fmt.Println("Invalid Node Number: Please input an integer value of number of nodes to test!")
		return 0, ""
	}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Would you like to print the details from each round: Yes(Y) or No(N)? ")
	_, err2 := fmt.Scanf("%s", &printResults)
	if err2 != nil || isValidPrinterSetting(printResults) {
		fmt.Println("Invalid Print Setting: Please enter either 'Y' or 'N' next time!")
		return 0, ""
	}
	return desiredNodes, printResults
}

// isValidPrinterSetting returns whether the user's choice of verbose output can be accepted by the program
func isValidPrinterSetting(printerSetting string) bool {
	if printerSetting != "Y" && printerSetting != "N" && printerSetting != "y" && printerSetting != "n" {
		return true
	}
	return false
}

// getAnalysis displays the number of rounds the gossip took to spread in every node if the user chooses verbose output
func getAnalysis(roundCount int, nodeCount int) {
	if printResults == "Y" || printResults == "y" {
		fmt.Println("------------------------------------------------------")
	}
	if printResults == "Y" || printResults == "y" {
		fmt.Printf("Well, it only took %d rounds to finish %d nodes.\n", roundCount, nodeCount)
	}
}
