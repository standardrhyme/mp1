package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func initiatePush(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	push(wg)
}

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
