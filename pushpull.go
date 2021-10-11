package main

import (
	"fmt"
	"sync"
)

/*
	@input wg //A pointer to the wait group
		   print //The user's choice of verbose output
	initiatePushPull displays the current number of round and calls the function that runs the push/pull gossip
*/
func initiatePushPull(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	pushPull(wg)
}

/*
	@input wg //A pointer to the wait group
		   print //The user's choice of verbose output
		   switchToPull //A boolean value to keep track of whether at least half of the nodes are infected
	@output complete //True if all nodes are infected, false otherwise
			switchToPull //True if more than half of the nodes are infected, false otherwise
	initiatePushPullSwitch constitutes one round of push/pull switch gossip. The function calls push gossip if less than
	half of the nodes are infected. It calls pull gossip otherwise.
*/
func initiatePushPullSwitch(wg *sync.WaitGroup, print string, switchToPull bool) (bool, bool) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	if !switchToPull {
		push(wg)
		complete, completeCount := completionCheck()
		if complete {
			return true, switchToPull
		}
		if completeCount*2 >= nodeCount {
			switchToPull = true
			if print == "Y" || print == "y" {
				fmt.Println("Switching to pull gossip.")
			}
		}
	} else {
		pull(wg)
		complete, _ := completionCheck()
		if complete {
			return true, switchToPull
		}
	}
	return false, switchToPull
}

/*
	@input wg //A pointer to the wait group
	pushPull constitutes one round of push/pull gossip. All nodes go through the following phases in order: push
	infection, push update, pull update, pull infection, clear channel.
*/
func pushPull(wg *sync.WaitGroup) {
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating push infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushInfect(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating push update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating pull update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating pull infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Clearing all channels.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go clearChannel(wg, &nodes[i])
	}
	wg.Wait()
}
