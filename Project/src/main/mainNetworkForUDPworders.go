package main

import (
	def "definitions"
	"fmt"
	"network"
	//"os"
	"time"
)

func main() {
	NetworkToQueueOrderChannel := make(chan def.ElevatorOrder, 10)
	NetworkToQueueAliveChannel := make(chan def.ElevatorAliveMessage, 10)
	QueueToNetworkOrderChannel := make(chan def.ElevatorOrder, 10)
	QueueToNetworkAliveChannel := make(chan def.ElevatorAliveMessage, 10)
	tick := time.Tick(5000 * time.Millisecond)
	tick2 := time.Tick(2100 * time.Millisecond)
	go network.InitUDP(NetworkToQueueOrderChannel, QueueToNetworkOrderChannel, NetworkToQueueAliveChannel, QueueToNetworkAliveChannel)

	var newOrder def.ElevatorOrder
	newOrder.Floor = 1
	newOrder.Btn = 1
	newOrder.ElevatorId = network.GetElevatorId()
	newOrder.Ack = false
	var newOrder2 def.ElevatorOrder
	newOrder2.Floor = 2
	newOrder2.Btn = 0
	newOrder2.ElevatorId = network.GetElevatorId()
	go PrintOrder(NetworkToQueueOrderChannel, QueueToNetworkOrderChannel)
	for {
		select {
		case <-tick:
			QueueToNetworkOrderChannel <- newOrder
			QueueToNetworkOrderChannel <- newOrder2
		case <-tick2:
			QueueToNetworkOrderChannel <- newOrder2
		}

	}

}

func PrintOrder(NetworkToQueueOrderChannel chan def.ElevatorOrder, QueueToNetworkOrderChannel chan def.ElevatorOrder) {
	for {
		orderToPrint := <-NetworkToQueueOrderChannel
		fmt.Println("NEW order recieved! Fantastic news\n")
		fmt.Println("The new order has ID %s \n", orderToPrint.ElevatorId)
		fmt.Println("Floor; %i \t Button; %i \t Ack: %t \n_________\n", orderToPrint.Floor, orderToPrint.Btn, orderToPrint.Ack)
		if orderToPrint.Ack == false {
			fmt.Println("The recieved order was appended to some kind of queue\nand is being sent to network with an ACK")
			orderToPrint.Ack = true
			QueueToNetworkOrderChannel <- orderToPrint
		}

	}

}
