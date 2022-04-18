package main

import (
	b "Blockchain"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"strconv"
	person "user_class"
)

var Users []person.User

var port int

var mining_port int

var blockchain *b.Block

func PrintUsers() {
	fmt.Println("Users in the blockchain")
	fmt.Printf("%-10s\t", "Name")
	fmt.Printf("%-10s\t", "Cash")
	fmt.Printf("%-15s\t", "Effort Points")
	fmt.Printf("%-15s\n", "Reward Points")
	for i := 0; i < len(Users); i++ {
		fmt.Printf("%-10s\t", Users[i].Name)
		fmt.Printf("%-10f\t", Users[i].Cash)
		fmt.Printf("%-10f\t", Users[i].EffortPoint)
		fmt.Printf("%-10f\n", Users[i].RewardPoint)
	}
}

func do_transaction() {
	for {
		fmt.Println("Who do you want to send coins to?")

		for i := 0; i < len(Users); i++ {
			if Users[i].Name != "Satoshi" {
				fmt.Println(strconv.Itoa(i) + ". " + Users[i].Name)
			}
		}

		var index int

		fmt.Scanf("%d\n", &index)

		fmt.Println("Enter the number of coins you want to send:")

		var coins int

		fmt.Scanf("%d\n", &coins)

		transaction_string := "Send " + strconv.Itoa(coins) + " coins to " + Users[index].Name

		fmt.Println("Send " + strconv.Itoa(coins) + " coins to " + Users[index].Name)

		var miner_index int

		max_priority := 0.0
		for i := 0; i < len(Users); i++ {
			priority := Users[i].EffortPoint / Users[i].RewardPoint

			if priority > max_priority {
				max_priority = priority
				miner_index = i
			}
		}

		conn1, err1 := net.Dial("tcp", ":"+strconv.Itoa(Users[miner_index].PortNumber))
		if err1 != nil {
			fmt.Println(err1)
		}

		handleConnectionMining(conn1, transaction_string, 0, index, miner_index, coins)
	}
}

func handleConnection1(c net.Conn) {
	gobEncoder := gob.NewEncoder(c)
	err := gobEncoder.Encode("INFO")
	err = gobEncoder.Encode(blockchain)
	if err != nil {
		log.Println(err)
	}

	fmt.Println()

	//fmt.Printf("%v\n", Users)

	err = gobEncoder.Encode(Users)

	if err != nil {
		log.Println(err)
	}
}

func handleConnectionMining(conn net.Conn, transaction_string string, sender_index int, receiver_index int, miner_index int, coins int) {
	enc := gob.NewEncoder(conn)
	err := enc.Encode("MINE")
	if err != nil {
		fmt.Println(err)
	}

	err = enc.Encode(transaction_string)
	if err != nil {
		fmt.Println(err)
	}

	err = enc.Encode(sender_index)
	if err != nil {
		fmt.Println(err)
	}
	err = enc.Encode(receiver_index)
	if err != nil {
		fmt.Println(err)
	}
	err = enc.Encode(miner_index)
	if err != nil {
		fmt.Println(err)
	}
	err = enc.Encode(coins)
	if err != nil {
		fmt.Println(err)
	}
}

func handleConnection(c net.Conn) {
	//var blockchain *a1.Block

	var name string
	port = port + 1
	mining_port = mining_port + 1
	dec := gob.NewDecoder(c)
	err := dec.Decode(&name)
	fmt.Println(name + " has connected on port " + strconv.Itoa(port))

	newUser := person.User{Name: name, Cash: 0, PortNumber: port, EffortPoint: 10.0, RewardPoint: 10.0}
	Users = append(Users, newUser)
	gobEncoder := gob.NewEncoder(c)
	err = gobEncoder.Encode(newUser.PortNumber)
	if err != nil {
		log.Println(err)
	}

	return
}

func Receive(conn net.Conn) {

	dec := gob.NewDecoder(conn)

	err := dec.Decode(&blockchain)

	if err != nil {
		fmt.Println(err)
	}

	err = dec.Decode(&Users)

	fmt.Println(b.CountBlocks(blockchain))

	if b.CountBlocks(blockchain)%5 == 0 && b.CountBlocks(blockchain) > 0 {
		for i := 0; i < len(Users); i++ {
			Users[i].EffortPoint = Users[i].EffortPoint * 0.9
			Users[i].RewardPoint = Users[i].RewardPoint * 0.9
		}
	}

	fmt.Println()

	b.ListBlocks(blockchain)

	PrintUsers()
}

func Mine(conn net.Conn) {

	var transaction_string string
	var sender_index int
	var receiver_index int
	var miner_index int
	var coins int

	dec := gob.NewDecoder(conn)
	err := dec.Decode(&transaction_string)
	if err != nil {
		fmt.Println(err)
	}
	err = dec.Decode(&sender_index)
	err = dec.Decode(&receiver_index)
	err = dec.Decode(&miner_index)
	err = dec.Decode(&coins)

	if b.VerifyChain(blockchain) {

		if Users[sender_index].Cash > float64(coins) {

			blockchain = b.InsertBlock(Users[receiver_index].Name+" "+strconv.Itoa(coins), Users[miner_index].Name+" "+"20", blockchain)

			Users[receiver_index].Cash += float64(coins)

			Users[sender_index].Cash -= float64(coins)

			Users[miner_index].Cash += 20

			for i := 0; i < len(Users); i++ {
				Users[i].EffortPoint += 15
			}

			Users[miner_index].RewardPoint += 10

			Users[sender_index].EffortPoint += 10

			for i := 0; i < len(Users); i++ {
				conn, err := net.Dial("tcp", ":"+strconv.Itoa(Users[i].PortNumber))
				if err != nil {
					fmt.Println(err)
				}

				handleConnection1(conn)
			}
		}
	}
}

func main() {
	blockchain = b.InsertBlock("Satoshi 100", "", nil)
	satoshi := person.User{Name: "Satoshi", Cash: 100, PortNumber: 6000, EffortPoint: 10.0, RewardPoint: 10.0}
	Users = append(Users, satoshi)
	max_users := 4
	mining_port = 9000
	port = 6000
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < max_users; i++ {
		conn, err := ln.Accept()
		blockchain = b.InsertBlock("Satoshi 20", "", blockchain)
		Users[0].Cash += 20
		if err != nil {
			log.Println(err)
		}
		go handleConnection(conn)
	}

	fmt.Println("Max users have been connected")

	for i := 0; i < len(Users); i++ {
		conn1, err1 := net.Dial("tcp", ":"+strconv.Itoa(Users[i].PortNumber))
		if err1 != nil {
			log.Println(err1)
		}
		go handleConnection1(conn1)

	}

	b.ListBlocks(blockchain)

	go do_transaction()

	for {
		var Action string
		conn1, err1 := ln.Accept()
		if err1 != nil {
			log.Println(err1)
		}
		dec := gob.NewDecoder(conn1)
		err := dec.Decode(&Action)
		if err != nil {
			fmt.Println(err)
		}

		if Action == "MINE" {
			go Mine(conn1)
		} else if Action == "INFO" {
			go Receive(conn1)
		}
	}
}
