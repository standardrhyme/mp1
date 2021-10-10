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
var desirednodes int
var desirednodesresults []int

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
	wg := &sync.WaitGroup{}
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Printf("Welcome! Which method of gossip would you like to implement: Push (PSH), Pull (PLL), or Push/Pull(PP)? ")
	fmt.Scanf("%s", &mode)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Printf("For up to how many nodes would you like to test the convergence speed? ")
	fmt.Scanf("%d", &desirednodes)
	if mode == "PSH" {
		isPushing = true
		isPulling = false
		initiate()
		pushgossip(desirednodes)
	} else if mode == "PLL" {
		isPushing = false
		isPulling = true
		initiate()
		pullgossip(desirednodes)
	} else {
		isPushing = true
		isPulling = true
		initiate()
		pushandpullgossip(desirednodes)
	}

	desirednodesresults = append(desirednodesresults, 0)
	for i := 1; i <= desirednodes; i++ {
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
			fmt.Println("------------------------------------------------------")
			fmt.Printf("Well, it only took %d rounds to finish %d nodes.\n", roundCount, nodeCount)
			desirednodesresults = append(desirednodesresults, roundCount)
		}
	}

	Plot()

}


package main

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

//Plot is a function that plots the number of nodes versus the convergence speed
func Plot() {
	var keys []int
	var values []opts.ScatterData

	for keyValue := 1; keyValue < len(desirednodesresults); keyValue++ {
		keys = append(keys, keyValue)
		values = append(values, opts.ScatterData{Value: desirednodesresults[keyValue]})
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(
			opts.Title{
				Title: "Nodes vs. Convergence Time",
			},
		),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "# of Nodes",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Conv. Time",
		}),
	)

	// Put data into instance
	scatter.SetXAxis(keys)
	scatter.AddSeries("Category A", values)
	scatter.SetSeriesOptions(charts.WithLabelOpts(
		opts.Label{
			Show:     true,
			Position: "right",
		}),
	)
	f, _ := os.Create("nodesvsconvergencetime.html")
	scatter.Render(f)

}
