package main

import (
	"fmt"
)

func getMode() string {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Welcome! Which method of gossip would you like to implement: " +
		"Push(1), Pull(2), or Push/Pull Original(3) or Push/Pull Switch(4)? Please enter the number code as indicated.")
	fmt.Println("If you would like to quit, please enter 'q'. ")
	_, err := fmt.Scanf("%s", &mode)
	if err != nil || isInvalidMode(mode) {
		fmt.Println("Invalid Code: Please select a valid gossip protocol next time!")
		return "QUIT"
	} else if mode == "q" || mode == "Q" {
		fmt.Println("You have quit the program. Thank you and goodbye!")
		return "QUIT"
	}
	return mode
}

func isInvalidMode(mode string) bool {
	if mode != "1" && mode != "2" && mode != "3" && mode != "4" && mode != "q" && mode != "Q" {
		return true
	}
	return false
}

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
	if err2 != nil || isInvalidPrinterSetting(printResults) {
		fmt.Println("Invalid Print Setting: Please enter either 'Y' or 'N' next time!")
		return 0, ""
	}
	return desiredNodes, printResults
}

func isInvalidPrinterSetting(printerSetting string) bool {
	if printerSetting != "Y" && printerSetting != "N" {
		return true
	}
	return false
}

func getAnalysis(roundCount int, nodeCount int) {
	if printResults == "Y" || printResults == "y" {
		fmt.Println("------------------------------------------------------")
	}
	if printResults == "Y" || printResults == "y" {
		fmt.Printf("Well, it only took %d rounds to finish %d nodes.\n", roundCount, nodeCount)
	}
}
