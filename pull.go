package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func initiatePull(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	pull(wg)
}

func pullInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if !node.infected {
		target := rand.Intn(nodeCount)
		select {
		case msg, ok := <-*nodes[target].channel:
			if ok {
				if print == "Y" || print == "y" {
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

func pullUpdate(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		for len(*(*node).channel) < nodeCount {
			*(*node).channel <- node.infected
		}
	}
}

func pull(wg *sync.WaitGroup) {
	if print == "Y" || print == "y" {
		fmt.Println("Initiating update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" || print == "y" {
		fmt.Println("Initiating infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
}
