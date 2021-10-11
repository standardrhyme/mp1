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
		for mode == "4" {
			roundCount++
			complete := initiatePushPullSwitch(wg, printResults, false)
			if complete {
				break
			}
		}
		desiredNodesResults = append(desiredNodesResults, roundCount)
	}
	return roundCount, nodeCount
}

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
