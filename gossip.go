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
var mode string
var isPulling bool
var isPushing bool
var pullSwitch bool
var desiredNodes int
var desirednodesresults []int
var print string

type Node struct {
	infected bool
	channel  *chan bool
}

func pushInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		//rand.Seed(time.Now().UnixNano())
		target := rand.Intn(nodeCount)
		if print == "Y" {
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

func pullInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if !node.infected {
		rand.Seed(time.Now().UnixNano())
		target := rand.Intn(nodeCount)
		select {
		case msg, ok := <-*nodes[target].channel:
			if ok {
				if print == "Y" {
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

func clearChannel(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	for len(*node.channel) > 0 {
		b := <-*node.channel
		node.infected = node.infected || b
	}
}

func push(wg *sync.WaitGroup) {
	if print == "Y" {
		fmt.Println("Initiating infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushInfect(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Initiating update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Clearing all channels.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go clearChannel(wg, &nodes[i])
	}
	wg.Wait()
}

func pull(wg *sync.WaitGroup) {
	if print == "Y" {
		fmt.Println("Initiating update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Initiating infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
}

func pushPull(wg *sync.WaitGroup) {
	if print == "Y" {
		fmt.Println("Initiating push infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Initiating push update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushInfect(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Initiating pull update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Initiating pull infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" {
		fmt.Println("Clearing all channels.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go clearChannel(wg, &nodes[i])
	}
	wg.Wait()
}

func completionCheck() (bool, int) {
	count := 0
	if print == "Y" {
		fmt.Println("Completion check.")
	}
	complete := true
	for i := 0; i < nodeCount; i++ {
		if print == "Y" {
			fmt.Printf("Node %d: %t\n", i, nodes[i].infected)
		}
		complete = nodes[i].infected && complete
		if nodes[i].infected {
			count++
		}
	}
	return complete, count
}

func main() {
	wg := &sync.WaitGroup{}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("Welcome! Which method of gossip would you like to implement: Push (PSH), Pull (PLL), or Push/Pull Original(PPO) or Push/Pull Switch (PPS)? Please enter the 3 character code as indicated.")
	fmt.Scanf("%s", &mode)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Println("For up to how many nodes would you like to test the convergence speed?")
	fmt.Scanf("%d", &desiredNodes)
	fmt.Println("Would you like to print the details from each round: Yes(Y) or No(N)? ")
	fmt.Scanf("%s", &print)
	desirednodesresults = append(desirednodesresults, 0)
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
		for mode == "PSH" {
			roundCount++
			if print == "Y" {
				fmt.Println("------------------------------------------------------")
				fmt.Printf("Round %d:\n", roundCount)
				fmt.Println("------------------------------------------------------")
			}
			push(wg)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//pull
		for mode == "PLL" {
			roundCount++
			if print == "Y" {
				fmt.Println("------------------------------------------------------")
				fmt.Printf("Round %d:\n", roundCount)
				fmt.Println("------------------------------------------------------")
			}
			pull(wg)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//push&pull
		for mode == "PPO" {
			roundCount++
			if print == "Y" {
				fmt.Println("------------------------------------------------------")
				fmt.Printf("Round %d:\n", roundCount)
				fmt.Println("------------------------------------------------------")
			}
			pushPull(wg)
			complete, _ := completionCheck()
			if complete {
				break
			}
		}
		//pull switch
		switchToPull := false
		for mode == "PPS" {
			roundCount++
			if print == "Y" {
				fmt.Println("------------------------------------------------------")
				fmt.Printf("Round %d:\n", roundCount)
				fmt.Println("------------------------------------------------------")
			}
			if !switchToPull {
				push(wg)
				complete, completeCount := completionCheck()
				if complete {
					break
				}
				if completeCount*2 >= nodeCount {
					switchToPull = true
					fmt.Println("Switching to pull gossip.")
				}
			} else {
				pull(wg)
				complete, _ := completionCheck()
				if complete {
					break
				}
			}
		}
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Well, it only took %d rounds to finish %d nodes.\n", roundCount, nodeCount)
		desirednodesresults = append(desirednodesresults, roundCount)
	}
	Plot(mode)

}
