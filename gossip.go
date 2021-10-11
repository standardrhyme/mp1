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
var printResults string

//Node is the entity that stores and spread gossip
type Node struct {
	infected bool
	channel  *chan bool
}

/*
	@input wg //A pointer to the wait group
	       node //A pointer to the node whose channel is to be cleared
	clearChannel concurrently resets the excessive values stored in the node's go channel
*/
func clearChannel(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	for len(*node.channel) > 0 {
		b := <-*node.channel
		node.infected = node.infected || b
	}
}

/*
	@output complete //A boolean value that indicates whether all nodes have been infected
			count //An integer showing the number of nodes that have been infected
	completionCheck check if the spread of the gossip has been complete
*/
func completionCheck() (bool, int) {
	count := 0
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Completion check.")
	}
	complete := true
	for i := 0; i < nodeCount; i++ {
		if printResults == "Y" || printResults == "y" {
			fmt.Printf("Node %d: %t\n", i, nodes[i].infected)
		}
		complete = nodes[i].infected && complete
		if nodes[i].infected {
			count++
		}
	}
	return complete, count
}

/*
	@input mode //The type of gossip chosen to be performed in this program run by the user
		   desiredNodes //The maximum number of nodes that will be tested in this run
		   printResults //Whether the user chooses to display the verbose logging of each round as output
		   wg //A pointer to the wait group
	initiateGossip processes all the user inputs and run the wanted gossip algorithm in accordance
*/
func initiateGossip(mode string, desiredNodes int, printResults string, wg *sync.WaitGroup) (int, int) {
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
			initiatePush(wg, printResults)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//pull
		for mode == "2" {
			roundCount++
			initiatePull(wg, printResults)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//push&pull
		for mode == "3" {
			roundCount++
			initiatePushPull(wg, printResults)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//push&pull--switch
		switchToPull := false
		for mode == "4" {
			roundCount++
			complete := false
			complete, switchToPull = initiatePushPullSwitch(wg, printResults, switchToPull)
			if complete {
				break
			}
		}
		desiredNodesResults = append(desiredNodesResults, roundCount)
	}
	return roundCount, nodeCount
}

// main function
func main() {
	wg := &sync.WaitGroup{}
	mode := getMode()
	if mode != "QUIT" {
		desiredNodes, printResults := getSettings()
		if desiredNodes != 0 {
			roundCount, nodeCount := initiateGossip(mode, desiredNodes, printResults, wg)
			getAnalysis(roundCount, nodeCount)
			Plot(mode)
		}
	}
}
