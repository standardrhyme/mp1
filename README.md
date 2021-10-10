# MP1
MP1 uses `Go-channels` and `Go-routines` to simulate three major types of Gossip protocols in a synchronous network.
These types are push-based gossip, pull-based gossip, and push-pull-based gossip. Two types of push-pull-based gossip are presented, making for a total of four possible gossip protocol options. 

## Input and Output 
Input is user-specified number of nodes in the system, and the output is the infection status of each node per round leading up to a fully infected network.

## Implementations of Push-Pull Gossip
#### Push Gossip (User Input Option 1)
To begin, each node will have a boolean false value as well as an associated channel. For example, node 1 will have an infected value of false and an associated channel 1. Rounds will then begin. Each node will run in its own goRoutine. PushInfect will begin, wherein each infected node sends a true value to the channel of a random node. The initialized wait group will then wait for all possible infections to occur. Then PushUpdate will begin, wherein each node will check their channel to see if they have received a true value. If they have received a true value, they will update their infected state to become infected. The wait group will wait for all possible updates to occur. Each node will then clear their channel to prepare for the next round. If all nodes have been infected, the loop will break and the number of rounds will be recorded and plotted.

#### Pull Gosip (User Input Option 2)
To begin, each ndoe will have a boolean false value as well as an associated channel (like Push Gossip above). Rounds will begin. Each node will run its own goRoutine. PullUpdate will begin, wherein each infected node sends a true value to its own channel. The initialized wait group will then wait for all nodes to update their own channels. Then PullInfect will begin, wherein each node that is not infected, will pull from a random nodes channel. If the random node's channel had a true value, the node that pulled will then update itself to show infected. Again, the wait group will wait for all possible infections to occur. Each node will then clear their channel to prepare for the next round. If all nodes have been infected, the loop will break and the number of rounds will be recorded and plotted.

#### Push/Pull Gossip (User Input Option 3)
Push/Pull gossip will begin the same as the other modes; with the initialization of nodes. Rounds will begin. PushInfect will first be called, working the same as it did in the Push Gossip. Then also will PushUpdate be called. Before the round ends, PullUpdate will be called, followed by PullInfect, each working the same as they did in the individual push and pull modes. 

#### Push/Pull Switch Gossip (User Input Option 4)
Push/Pull Switch gossip will begin the same as the other modes; with the initialization of nodes. Rounds will begin. While the number of susceptible nodes is more than half the number of total nodes, the push gossip will be implemented. Once the number of susceptible nodes reaches less than half, the pull gossip will be implemented. The switch is made when around half the number of nodes are susceptible because of the probability of a node being randomly chosen. To acheive a smaller convergence time, the probability of a susceptible node being infected in a round should be higher than not being infected. When the number of susceptible nodes is greater than the number of infected, there is a greater probability that a susceptible node will turn infected. This is because the already infected nodes will have a higher probability of choosing a susceptible node for I(t) (the number of infected nodes at time t) < S(t) (the number of susceptible nodes at time t) and thus S(t)/(I(t) + S(t)) will be greater than I(t)/(I(t) + S(t)). On the other hand, when the number of susceptible nodes is smaller than the number of infected, there is a greater probability that a susceptible node will pull from an infected than an infected node pushing to a susceptible. In other terms, I(t)/(I(t) + S(t)). 

## How to Run

### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/standardrhyme/mp1`.

### Step 2: Initialize gossip protocol 
Change the current directory to be within the recently cloned folder. Start the gossip protocol with `go run mp1`.

##### If an error of the following form (plot.go:7:2: cannot find package "github.com/go-echarts/go-echarts/v2/charts" in any of: /usr/local/Cellar/go/1.17/...) is triggered, run `export GO111MODULE=on`.

### Step 3: Interact with Command Line
Enter an integer value of the number of nodes you want to infect.  
Press ENTER to begin gossip protocol. 

If you wish to quit the program, enter `q` or `Q`.

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

## Exit Codes 
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References 
- The plotting function `plot` is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
