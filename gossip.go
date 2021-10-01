package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var nodeCount int
var nodes []Node

type Node struct {
	infected bool
	channel *chan bool
}

func pushInfect(wg *sync.WaitGroup, node *Node)  {
	defer wg.Done()
	if node.infected {
		target := rand.Intn(nodeCount)
		fmt.Printf("Node %d is being infected.\n", target)
		*nodes[target].channel <- true
	}
}

func pushUpdate(wg *sync.WaitGroup, node *Node)  {
	defer wg.Done()
	select {
	case msg, ok := <- *node.channel:
		if ok {
			node.infected = msg
		} else {
			fmt.Println("Something wrong with the message")
			break
		}
	default:
		break
	}
	for len(*node.channel) > 0 {
		<- *node.channel
	}
}

func main ()  {
	wg := &sync.WaitGroup{}
	nodeCount = 100
	nodes = make([]Node, nodeCount)
	channels := make([]chan bool, nodeCount)
	for i := 0; i < nodeCount; i++ {
		channels[i] = make(chan bool, nodeCount)
		nodes[i] = Node{false, &(channels[i])}
	}
	nodes[0].infected = true
	roundCount := 0
	for true {
		roundCount ++
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
		fmt.Println("Initiating infection phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pushInfect(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Initiating update phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pushUpdate(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Completion check: ")
		complete := true
		for i := 0; i < nodeCount; i++ {
			fmt.Printf("%d:%t\n", i, nodes[i].infected)
			complete = nodes[i].infected && complete
		}
		if complete {
			break
		}
	}
	fmt.Println("------------------------------------------------------")
	fmt.Printf("Well, it only took %d rounds to finish %d nodes.", roundCount, nodeCount)
}