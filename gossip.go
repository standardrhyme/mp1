package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var nodeCount int
var roundCount int
var nodes []Node
var isPulling bool
var isPushing bool

type Node struct {
	infected bool
	channel  *chan bool
}

func pushInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		rand.Seed(time.Now().UnixNano())
		target := rand.Intn(nodeCount)
		fmt.Printf("Node %d is being infected.\n", target)
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
	for len(*node.channel) > 0 {
		<-*node.channel
	}
}

func pullInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if !node.infected {
		rand.Seed(time.Now().UnixNano())
		target := rand.Intn(nodeCount)
		select {
		case msg, ok := <-*nodes[target].channel:
			if ok {
				fmt.Printf("A node is being infected by node %d. \n", target)
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

func clearChannel(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	for len(*node.channel) > 0 {
		<-*node.channel
	}
}

func main() {
	isPushing = true
	isPulling = true
	wg := &sync.WaitGroup{}
	nodeCount = 10
	nodes = make([]Node, nodeCount)
	channels := make([]chan bool, nodeCount)
	for i := 0; i < nodeCount; i++ {
		channels[i] = make(chan bool, nodeCount)
		nodes[i] = Node{false, &(channels[i])}
	}
	nodes[0].infected = true
	roundCount = 0
	//push
	for isPushing && !isPulling {
		roundCount++
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
	//pull
	for isPulling && !isPushing {
		roundCount++
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
		fmt.Println("Initiating update phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pullUpdate(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Initiating infection phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pullInfect(wg, &nodes[i])
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
	//push&pull
	for isPulling && isPushing {
		roundCount++
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
		fmt.Println("Initiating push infection phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pushUpdate(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Initiating push update phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pushInfect(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Initiating pull update phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pullUpdate(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Initiating pull infection phase:")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go pullInfect(wg, &nodes[i])
		}
		wg.Wait()
		fmt.Println("Clearing all channels.")
		for i := 0; i < nodeCount; i++ {
			wg.Add(1)
			go clearChannel(wg, &nodes[i])
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
