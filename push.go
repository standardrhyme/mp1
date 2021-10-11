package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*
	@input wg //A pointer to the wait group
		   print //The user's choice of verbose output
	initiatePush displays the current number of round and calls the function that runs the push gossip
*/
func initiatePush(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	push(wg)
}

/*
	@input wg //A pointer to the wait group
		   node //The node on which the action of push infection is to be performed
	If the input node is infected, pushInfect chooses another node's channel and push the gossip
*/
func pushInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		target := rand.Intn(nodeCount)
		if printResults == "Y" || printResults == "y" {
			fmt.Printf("Node %d is being infected.\n", target)
		}
		*nodes[target].channel <- node.infected
	}
}

/*
	@input wg //A pointer to the wait group
		   node //The node on which the action of push infection is to be performed
	If the input node is susceptible, pushUpdate checks its channel for any possible incoming gossip
*/
func pushUpdate(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	select {
	case msg, ok := <-*node.channel:
		if ok {
			node.infected = msg
		} else {
			fmt.Println("Channel closed for some reason.")
			break
		}
	default:
		break
	}
}

/*
	@input wg //A pointer to the wait group
	push constitutes one round of push gossip, where all nodes go through the infection phase and then the update phase
	before clearing the channel
*/
func push(wg *sync.WaitGroup) {
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushInfect(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushUpdate(wg, &nodes[i])
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
