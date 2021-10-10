# MP1
MP1 uses `Go-channels` and `Go-routines` to simulate three major types of Gossip protocols in a synchronous network.
These types are push-based gossip, pull-based gossip, and push-pull-based gossip. Two types of push-pull-based gossip are presented, making for a total of four possible gossip protocol options. 

## Input and Output 
Input is user-specified number of nodes in the system, and the output is the infection status of each node per round leading up to a fully infected network.

## Implementations of Push-Pull Gossip


## How to Run

### Step 1: Clone Git Repository
Clone the following git repository with `git clone https://github.com/standardrhyme/mp1`.

### Step 2: Initialize gossip protocol 
Change the current directory to be within the recently cloned folder. Start the gossip protocol with `go run gossip.go plot.go`

##### If an error of the following form is triggered:

##### If the error is not solved, install the following dependencies with the following: 

### Step 3: Interact with Command Line
Enter an integer value of the number of nodes you want to infect.  
Press ENTER to begin gossip protocol. 

If you wish to quit the program, enter `q`.

## Screenshots

1. Command Line Interface - Valid User Input

2. Command Line Interface - User Quit Program

3. Output

## Workflows

General Workflow

Push Gossip Overview

Pull Gossip Overview

Push-Pull Gossip Overview

## Custom Data Structures

## Exit Codes 
- `0`: Successful
- `1`: Incorrect command line input format
- `2`: External package function error

## References 
- The plotting function `plot` is a modified version of sample code from [Go E-Charts Examples](https://github.com/go-echarts/examples/blob/master/examples/scatter.go "Go E-Charts Examples").
