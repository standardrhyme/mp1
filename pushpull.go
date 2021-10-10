package main

import (
	"fmt"
	"sync"
)

func initiatePushPull(wg *sync.WaitGroup, print string) {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	pushPull(wg)
}

func initiatePushPullSwitch(wg *sync.WaitGroup, print string, switchToPull bool) bool {
	if print == "Y" || print == "y" {
		fmt.Println("------------------------------------------------------")
		fmt.Printf("Round %d:\n", roundCount)
		fmt.Println("------------------------------------------------------")
	}
	if !switchToPull {
		push(wg)
		complete, completeCount := completionCheck()
		if complete {
			return true
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
			return true
		}
	}
	return false
}

func pushPull(wg *sync.WaitGroup) {
	if print == "Y" || print == "y" {
		fmt.Println("Initiating push infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushInfect(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" || print == "y" {
		fmt.Println("Initiating push update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pushUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" || print == "y" {
		fmt.Println("Initiating pull update phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullUpdate(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" || print == "y" {
		fmt.Println("Initiating pull infection phase.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go pullInfect(wg, &nodes[i])
	}
	wg.Wait()
	if print == "Y" || print == "y" {
		fmt.Println("Clearing all channels.")
	}
	for i := 0; i < nodeCount; i++ {
		wg.Add(1)
		go clearChannel(wg, &nodes[i])
	}
	wg.Wait()
}
