package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math/rand"
	"os"
	"sync"
)

var nodeCount int
var roundCount int
var nodes []Node
var isPulling bool
var isPushing bool
var pullSwitch bool
var desiredNodes int
var desirednodesresults []int

type Node struct {
	infected bool
	channel  *chan bool
}

func pushInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if node.infected {
		//rand.Seed(time.Now().UnixNano())
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
}

func pullInfect(wg *sync.WaitGroup, node *Node) {
	defer wg.Done()
	if !node.infected {
		//rand.Seed(time.Now().UnixNano())
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
		b := <-*node.channel
		node.infected = node.infected || b
	}
}

func push(wg *sync.WaitGroup) {
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
	fmt.Println("Clearing all channels.")
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go clearChannel(wg, &nodes[i])
	}
	wg.Wait()
}

func pull(wg *sync.WaitGroup) {
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
}

func pushPull(wg *sync.WaitGroup) {
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
}

func completionCheck() (bool, int) {
	count := 0
	fmt.Println("Completion check: ")
	complete := true
	for i := 0; i < nodeCount; i++ {
		fmt.Printf("%d:%t\n", i, nodes[i].infected)
		complete = nodes[i].infected && complete
		if nodes[i].infected {
			count++
		}
	}
	return complete, count
}

func plot() {
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
	err := scatter.Render(f)
	if err != nil {
		return
	}

}

func main() {
	isPushing = false
	isPulling = false
	pullSwitch = true
	wg := &sync.WaitGroup{}
	// fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	// //Get the desired number of nodes from the user
	// fmt.Println("How many nodes would you like the map to contain? To quit the program, enter '-1'.")
	// fmt.Scanf("%d", &nodeCount)
	// //If the user inputs 0, ask again.
	// for nodeCount == 0 {
	// 	fmt.Println("How many nodes would you like the map to contain? To quit the program, enter '-1'.")
	// 	fmt.Scanf("%d", &nodeCount)
	// }
	fmt.Println("Up to how many nodes do you want to test?")
	_, err := fmt.Scanf("%d", &desiredNodes)
	if err != nil {
		return
	}
	desirednodesresults = append(desirednodesresults, 0)
	for i := 1; i <= desiredNodes; i++ {
		nodeCount = i
		//If the user quits the program.
		if nodeCount != -1 {
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
				push(wg)
				complete, _ := completionCheck()
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
				pull(wg)
				complete, _ := completionCheck()
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
				pushPull(wg)
				complete, _ := completionCheck()
				if complete {
					break
				}
			}
			//pull switch
			switchToPull := false
			for pullSwitch && !isPulling && !isPushing {
				roundCount++
				fmt.Println("------------------------------------------------------")
				fmt.Printf("Round %d:\n", roundCount)
				fmt.Println("------------------------------------------------------")
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
		} else {
			fmt.Println("You have exited the program by entering '-1'. Thank you and goodbye!")
		}
	}

	plot()

}
