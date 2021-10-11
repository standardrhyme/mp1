# MP1
MP1 uses `Go-channels` and `Go-routines` to simulate three major types of Gossip protocols in a synchronous network.
These types are push-based Gossip, pull-based Gossip, and push-pull-based Gossip. Two types of push-pull-based gossip are presented, making for a total of four possible Gossip protocol options. 

## Input and Output 

## Specifications of Gossip Protocols

#### Initialization of Nodes
Once the user inputs the type of Gossip protocol to run, the number of nodes to initialize, and whether to print the round results in the terminal, a call to `initiateGossip()` is made. This function first creates an array of Nodes, with a length determine by user-input number of nodes. Each node in the array is initialized with an `infected` value of `false`, representing its infection status, as well as its own personal `bool channel`, which will be used for communication with other nodes. For example, Node 1 will be initialized with an infected value of `false` and its own `channel`. Aftewards, the program executes the instructions specific to the user-chosen Gossip protocol:

#### Push Gossip 
During each round, each node runs three total `Go-routines` in lockstep. The first is `PushInfect()`: if a node is infected, then it randomly selects a peer node and sends a `true` value to this peer's `channel`. The second is `PushUpdate()`: each node checks its channel to see if it has received one or more `true` values from its peer nodes. If it has, then it updates its infected state to `true`, to reflect that it has been infected. The third is `clearChannel()`: each node clears its `bool channel` using the `clearChannel()` function, which iterates through the node's channel and re-assigns the node's infected status to these values. Because only `true` values are ever passed into the channel, there is a net-zero effect on the ultimate infected status. Each of these three `Go-routine` functions are performed in lockstep by all nodes, made possible by synchronizing their completion through the `sync.WaitGroup` package. Once all nodes have been infected, which is checked each round by a call to `completionCheck()`, the `push` loop in `initiateGossip()` breaks, and the number of rounds is recorded and passed to `Plot()` to be plotted. An `HTML` file is created that maps `n`des to Convergence Time.

#### Pull Gossip 
As in the `Push` protocol, each node runs three total `Go-routines` in lockstep. The first is `PullUpdate()`: if a node is infected, it sends a `true` value to its own channel. The second is `PullInfect()`: if a node is not infected, it selects a random peer node and makes a pull request from this peer's `bool channel`. If the peer node's channel has a true value, the pulling node then updates its infected status to `true`. The third is `clearChannel()`, which is executed as in the `Push` protocol to clear the channels for the next round. Once all nodes have been infected, which is confirmed by `completionCheck()`, the `pull` loop in `initiateGossip()` breaks, and the number of rounds is recorded and passed to `Plot()` to be plotted. An `HTML` file is created that maps `n`des to Convergence Time.

#### Push/Pull Gossip 
`Push/Pull` combines the functions of `Push` and `Pull` protocols, as the name suggests. `PushInfect()` is called first in `Go-routines` for by each node, followed by `PushUpdate()`. Then, `PullUpdate()` and `PullInfect()` are called, completing the round. Each of them are executed in lockstep as in the previous protocols, and the outputs are created in the same manner.  

#### Push/Pull Switch Gossip 
Push/Pull Switch Gossip will begin the same as the other modes; with the initialization of nodes. Rounds will begin. While the number of susceptible nodes is more than half the number of total nodes, the push Gossip will be implemented. Once the number of susceptible nodes reaches less than half, the pull Gossip will be implemented. The switch is made when around half the number of nodes are susceptible because of the probability of a node being randomly chosen. To acheive a smaller convergence time, the probability of a susceptible node being infected in a round should be higher than not being infected. When the number of susceptible nodes is greater than the number of infected, there is a greater probability that a susceptible node will turn infected. This is because the already infected nodes will have a higher probability of choosing a susceptible node for I(t) (the number of infected nodes at time t) < S(t) (the number of susceptible nodes at time t) and thus S(t)/(I(t) + S(t)) will be greater than I(t)/(I(t) + S(t)). On the other hand, when the number of susceptible nodes is smaller than the number of infected, there is a greater probability that a susceptible node will pull from an infected than an infected node pushing to a susceptible. In other terms, I(t)/(I(t) + S(t)). 

## How to Run

### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/standardrhyme/mp1`.

### Step 2: Initialize Gossip Protocol
Change the current directory into the recently cloned `mp1` folder. Start the Gossip protocol with `go run mp1`.

##### If an error of the following form (plot.go:7:2: cannot find package "github.com/go-echarts/go-echarts/v2/charts" in any of: /usr/local/Cellar/go/1.17/...) is triggered, run `export GO111MODULE=on`.

### Step 3: Interact with Command Line
A) Enter the Integer Code corresponding to the type of Gossip protocol you wish to implement, and press `ENTER`.
 - `1`: Push
 - `2`: Pull
 - `3`: Push/Pull Original
 - `4`: Push/Pull Switch

If you wish to quit the program, enter `q` or `Q`.

B) Next, enter a postive integer value of the number of nodes you want in your system, and press `ENTER`. 

C) Lastly, enter whether you wish to print out in your terminal the infection results of each Gossip round, and press `ENTER`.
- `Y`: Yes
- `N`: No

## Screenshots
1. Command Line Interface - Valid User Input

2. Command Line Interface - User Quit Program

3. Output

## Workflows

General Workflow

<img src="https://user-images.githubusercontent.com/60116121/136716646-9d5d557a-5a53-4d59-b5c1-fec69e3b77aa.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136716646-9d5d557a-5a53-4d59-b5c1-fec69e3b77aa.png" width="50%" height="50%" />

Push Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713597-1a680e8b-d028-4d11-8717-ea2ae3538882.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713597-1a680e8b-d028-4d11-8717-ea2ae3538882.png" width="50%" height="50%" />

Pull Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713589-4a5952c5-0a8b-4a84-99d4-5eabadfb3568.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713589-4a5952c5-0a8b-4a84-99d4-5eabadfb3568.png" width="50%" height="50%" />

Push-Pull Gossip Overview

<img src="https://user-images.githubusercontent.com/60116121/136713592-ef8767b3-e920-4b83-9a14-218b43423169.png" data-canonical-src="https://user-images.githubusercontent.com/60116121/136713592-ef8767b3-e920-4b83-9a14-218b43423169.png" width="50%" height="50%" />

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
