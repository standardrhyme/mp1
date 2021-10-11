package main

import (
	"fmt"
	"math/rand"
	"sync"
)

/*
	@input wg //A pointer to the wait group
		   print //The user's choice of verbose output
	initiatePull displays the current number of round and calls the function that runs the pull gossip
*/
func initiatePull(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	pull(wg)
}

/*
	@input wg //A pointer to the wait group
		   node //The node on which the action of pull infection is to be performed
	If the input node is susceptible, pullInfect chooses another node's channel and attempt to pull the gossip
*/
func pullInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if !node.infected {
		target := rand.Intn(nodeCount)
		select {
		case msg, ok := <-*nodes[target].channel:
			if ok {
				if printResults == "Y" || printResults == "y" {
					fmt.Printf("A node is being infected by node %d. \n", target)
				}
				node.infected = msg
			} else {
				fmt.Println("Channel closed for some reason.")
				break
			}
		default:
			break

		}
	}
}

/*
	@input wg //A pointer to the wait group
		   node //The node on which the action of pull update is to be performed
	If the input node is infected, pullUpdate fills its channel with gossip
*/
func pullUpdate(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		for len(*(*node).channel) < nodeCount {
			*(*node).channel <- node.infected
		}
	}
}

/*
	@input wg //A pointer to the wait group
	pull constitutes one round of pull gossip, where all nodes go through the update phase and then the infection phase
*/
func pull(wg *sync.WaitGroup) {
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if printResults == "Y" || printResults == "y" {
		fmt.Println("Initiating infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
}
