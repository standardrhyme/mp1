package main

import (
	"fmt"
	"sync"
)

var nodeCount int
var roundCount int
var nodes []Node
var mode string
var desiredNodes int
var desiredNodesResults []int
var print string

type Node struct {
	infected bool
	channel  *chan bool
}

func clearChannel(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	for len(*node.channel) > 0 {
		b := <-*node.channel
		node.infected = node.infected || b
	}
}

func completionCheck() (bool, int) {
	count := 0
	if print == "Y" || print == "y" {
		fmt.Println("Completion check.")
	}
	complete := true
	for i := 0; i < nodeCount; i++ {
		if print == "Y" || print == "y" {
			fmt.Printf("Node %d: %t\n", i, nodes[i].infected)
		}
		complete = nodes[i].infected && complete
		if nodes[i].infected {
			count++
		}
	}
	return complete, count
}

func getAnalysis(roundCount int, nodeCount int) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
	}
	if print == "Y" || print == "y" {
		fmt.Printf("Well, it only took %d rounds to finish %d nodes.\n", roundCount, nodeCount)
	}
}

func getSettings() (int, string) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("For up to how many nodes would you like to test the convergence speed?")
	fmt.Scanf("%d", &desiredNodes)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Would you like to print the details from each round: Yes(Y) or No(N)? ")
	fmt.Scanf("%s", &print)
	return desiredNodes, print
}

func initiateGossip(mode string, desiredNodes int, print string, wg *sync.WaitGroup) (int, int) {
	desiredNodesResults = append(desiredNodesResults, 0)

	for i := 1; i <= desiredNodes; i++ {
		nodeCount = i
		nodes = make([]Node, nodeCount)
		channels := make([]chan bool, nodeCount)
		for i := 0; i < nodeCount; i++ {
			channels[i] = make(chan bool, nodeCount)
			nodes[i] = Node{false, &(channels[i])}
		}
		nodes[0].infected = true
		roundCount = 0

		//push
		for mode == "1" {
			roundCount++
			initiatePush(wg, print)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//pull
		for mode == "2" {
			roundCount++
			initiatePull(wg, print)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//push&pull
		for mode == "3" {
			roundCount++
			initiatePushPull(wg, print)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//pull switch
		for mode == "4" {
			switchToPull := false
			roundCount++
			complete := initiatePushPullSwitch(wg, print, switchToPull)
			if complete {
				break
			}
		}
		desiredNodesResults = append(desiredNodesResults, roundCount)
	}
	return roundCount, nodeCount
}

func isInvalidMode(mode string) bool {
	if mode != "1" && mode != "2" && mode != "3" && mode != "4" && mode != "q" && mode != "Q" {
		return true
	}
	return false
}

func getMode() string {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Welcome! Which method of gossip would you like to implement: " +
		"Push(1), Pull(2), or Push/Pull Original(3) or Push/Pull Switch(4)? Please enter the number code as indicated.")
	fmt.Println("If you would like to quit, please enter 'q'. ")
	fmt.Scanf("%s", &mode)
	error := isInvalidMode(mode)
	if error {
		fmt.Println("Invalid Mode: Please select a valid gossip protocol next time!")
		return "QUIT"
	} else if mode == "q" || mode == "Q" {
		fmt.Println("You have quit the program. Thank you and goodbye!")
		return "QUIT"
	}
	return mode
}

func main() {
	wg := &sync.WaitGroup{}
	mode := getMode()
	if mode != "QUIT" {
		desiredNodes, print := getSettings()
		roundCount, nodeCount := initiateGossip(mode, desiredNodes, print, wg)
		getAnalysis(roundCount, nodeCount)
		Plot(mode)
	}
}
