# BlockchainSimulation

## Introduction

This project simulates a Blockchain where we can see how a blockchain works in a distributed environment and how a node mines a block on the basis of a consensus algorithm.

## Setup

Clone this repository inside your GOPATH in the src folder and then run

$ go build main/main.go
$ go build node/node.go

After that run the main file

$ ./main

While main file is running, add nodes in the distributed system with command

$ ./node name_of_node
  
A minimum of 4 nodes must be added before the distributed system is able to run. After 4 nodes are added, we can simulate transactions by sending coins between different nodes and see which nodes gets to mine the block.
  
## Consensus Algorithm
  
### Name
  
Proof of Effort (POE)
  
### Objective
  
To reward users who participate and put in more effort in the blockchain.
  
### Description
  
All users in the blockchain get effort points for being part of the blockchain, but users who perform transactions will get relatively more effort points than the rest of the users. The miner is selected based on their priority (ratio between effort points and reward points (EP/RP = Priority). The miner also receives reward points for mining a block which will eventually affect their priority. The EP and RP are given as follows:
User who performs transaction = +25 EP
Users who are part of blockchain but did not perform the transaction = +15 EP
Miner = +10 RP & +15 EP
To balance the priorities and the EPRP system in the long run we will decay every users’ EP and RP by 10% every 5 blocks mined, this will give a relatively fair chance to all in the blockchain and encourage regular participation in the blockchain.
