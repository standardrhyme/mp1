# MP1
MP1 uses `Go-channels` and `Go-routines` to simulate three major types of Gossip protocols in a synchronous network.
These types are push-based Gossip, pull-based Gossip, and push-pull-based Gossip. Two types of push-pull-based gossip are presented, making for a total of four possible Gossip protocol options. 

## Input and Output
A user will specify which Gossip protocol they would like to see tested. Additionally, the user will provide a value for the maximum number of nodes in a system to test. Finally, the user will select whether they would like to see the results for each round printed (Y or N). 

Depending on the user's answer to the last question, the program will print the results of each infection round. Also, if the user's current directory is writeable, an HTML file will be output. The HTML file will display a plot with the number of nodes tested in a system against the number of rounds it took to infect all the nodes in that system.

## Specifications of Gossip Protocols

#### Initialization of Nodes
Once the user inputs the type of Gossip protocol to run, the number of nodes to test until, and whether to print the round results in the terminal, a call to `initiateGossip()` is made. To begin, a variable i (the number of nodes in the currently tested system) is set to 1. The function will then begin a loop wherein it first creates an array (of length i) of Nodes. Each node in that array, except the first, is initialized with an `infected` value of `false`, representing its infection status, as well as its own personal `bool channel`, which will be used for communication with other nodes. The first node, however, is initialized with an `infected` value of `true`, as well as its own personal `bool channel`. Afterwards, the program executes the instructions specific to the user-chosen Gossip protocol until all the nodes in the system are infected. Once those results are recorded, the loop will be run again with i+1 nodes in a system. This loop is run until the variable i has reached the number of nodes the user chose to test until.

For example, if the user inputs 50 for the number of nodes, the following systems will be tested: a system with 1 node, a system with 2 nodes, a system with 3 nodes,..., a system with 50 nodes.

#### Push Gossip 
During each round, each node runs three total `Go-routines` in lockstep. The first is `PushInfect()`: if a node is infected, then it randomly selects a peer node and sends a `true` value to this peer's `channel`. The second is `PushUpdate()`: each node checks its channel to see if it has received one or more `true` values from its peer nodes. If it has, then it updates its infected state to `true`, to reflect that it has been infected. The third is `clearChannel()`: each node clears its `bool channel` using the `clearChannel()` function, which iterates through the node's channel and re-assigns the node's infected status to these values. Because only `true` values are ever passed into the channel, there is a net-zero effect on the ultimate infected status. Each of these three `Go-routine` functions are performed in lockstep by all nodes, made possible by synchronizing their completion through the `sync.WaitGroup` package. Once all nodes have been infected, which is checked each round by a call to `completionCheck()`, the `push` loop in `initiateGossip()` breaks, and the number of rounds is recorded and passed to `Plot()` to be plotted. An `HTML` file is created that maps number of nodes to number of rounds (convergence time).

#### Pull Gossip 
As in the `Push` protocol, each node runs three total `Go-routines` in lockstep. The first is `PullUpdate()`: if a node is infected, it sends as many `true` values as their are nodes to its own channel. The second is `PullInfect()`: if a node is not infected, it selects a random peer node and makes a pull request from this peer's `bool channel`. If the peer node's channel has a true value, the pulling node then updates its infected status to `true`. The third is `clearChannel()`, which is executed as in the `Push` protocol to clear the channels for the next round. Once all nodes have been infected, which is confirmed by `completionCheck()`, the `pull` loop in `initiateGossip()` breaks, and the number of rounds is recorded and passed to `Plot()` to be plotted. An `HTML` file is created that maps nodes to Convergence Time.

#### Push/Pull Gossip 
`Push/Pull` combines the functions of `Push` and `Pull` protocols, as the name suggests. `PushInfect()` is called first in `Go-routines` for by each node, followed by `PushUpdate()`. Then, `PullUpdate()` and `PullInfect()` are called, completing the round. Each of them are executed in lockstep as in the previous protocols, and the outputs are created in the same manner.  

#### Push/Pull Switch Gossip 
This version of `Push/Pull` is implemented in this MP to demonstrate a more efficient variation of the original `Push/Pull` protocol. 

`Push/Pull Switch Gossip` is identical to `Push/Pull`, except that instead of executing both `Push` and `Pull` commands every round, it acknowledges the different phases of a general Gossip protocol, and implements the commands that are optimal to the current phase. More specifically, while the number of susceptible nodes is more than half the number of total nodes, the `Push` protocol is implemented. Conversely, while the number of susceptible nodes is less than half, the `Pull` protocol is implemented. 

This timing of the switch from `Push` to `Pull` makes optimal use of random peer node selection. `Push` works best when random peer selection is more likely to select susceptible nodes. `Pull` works best when random peer selection is more likely to select infected nodes. So, when the number of susceptible nodes is greater than the number of infected, infected nodes are more likely to select and infect a susceptible node, rather than re-selecting an already-infected node--`Push` Gossip is used. On the other hand, as the infection spreads and the number of infected nodes is greater than the number of susceptible nodes, susceptible nodes are more likely to select and pull the gossip from an infected node, rather than selecting another susceptible node to no net-effect--`Pull` Gossip is used. 

## How to Run

#### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/standardrhyme/mp1`.

#### Step 2: Initialize Gossip Protocol
Change the current directory into the recently cloned `mp1` folder. Start the Gossip protocol with `go run mp1`. 

###### If an error of the following form (plot.go:7:2: cannot find package "github.com/go-echarts/go-echarts/v2/charts" in any of: /usr/local/Cellar/go/1.17/...) is triggered, run `export GO111MODULE=on`.

###### If an error of the following form (cannot find package "mp1") is triggered, start the gossip protocol with `go run .`

#### Step 3: Interact with Command Line
A) Enter the Integer Code corresponding to the type of Gossip protocol you wish to implement, and press `ENTER`.
 - `1`: Push
 - `2`: Pull
 - `3`: Push/Pull Original
 - `4`: Push/Pull Switch

If you wish to quit the program, enter `q` or `Q`.

B) Next, enter a postive integer value of the number of nodes you want your system to test, and press `ENTER`. 

C) Lastly, enter whether you wish to print out in your terminal the infection results of each Gossip round, and press `ENTER`.
- `Y`: Yes
- `N`: No

## Screenshots
#### Command Line Interface - Valid User Input
<img width="902" alt="Screen Shot 2021-10-11 at 11 03 26 AM" src="https://user-images.githubusercontent.com/60116121/136813285-c3236b87-dcef-45aa-9da1-50612ff464dd.png">


#### Command Line Interface - User Quit Program
<img width="896" alt="Screen Shot 2021-10-11 at 11 03 53 AM" src="https://user-images.githubusercontent.com/60116121/136813355-52072136-fa0f-4470-8a55-ce4b30001527.png">


#### Output
###### If the user indicates they would like the round results to be printed.
<img width="897" alt="Screen Shot 2021-10-11 at 11 05 13 AM" src="https://user-images.githubusercontent.com/60116121/136813557-a44f82aa-1062-4607-9e46-dec076cfc72a.png">


###### Nodes vs. Convergence Time Results 
<img width="919" alt="Screen Shot 2021-10-11 at 11 06 11 AM" src="https://user-images.githubusercontent.com/60116121/136813708-5720c6a2-ef26-4670-a850-6a6c5e749710.png">


## Workflows

#### General Workflow

<img src="https://user-images.githubusercontent.com/60116121/137044276-0d4285f7-f7c0-49c3-a198-5fff84fdba81.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/137044276-0d4285f7-f7c0-49c3-a198-5fff84fdba81.png" width="100%" height="100%" />

#### Push Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713597-1a680e8b-d028-4d11-8717-ea2ae3538882.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713597-1a680e8b-d028-4d11-8717-ea2ae3538882.png" width="75%" height="75%" />

#### Pull Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713589-4a5952c5-0a8b-4a84-99d4-5eabadfb3568.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713589-4a5952c5-0a8b-4a84-99d4-5eabadfb3568.png" width="75%" height="75%" />

#### Push-Pull Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713592-ef8767b3-e920-4b83-9a14-218b43423169.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713592-ef8767b3-e920-4b83-9a14-218b43423169.png" width="75%" height="75%" />

## Custom Data Structures
Node Struct in ```Go ```
```
type Node struct {
  infected bool
  channel *chan bool
}
```

## Exit Codes 
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References 
- The plotting function `plot` in `plot.go` is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
