package main

import (
	def "definitions"
	"fmt"
	"network"
	//"os"
	"time"
)

func main() {
	msgCh := def.MessageChan{
		Outgoing: 		make(chan def.Message, 10),
		Incoming: 		make(chan def.Message, 10),
		CostReply: 		make(chan def.Message, 10),
		NumOnline: 		make(chan def.NumOnline),
	}
	
	go network.InitUDP(msgCh.Incoming,msgCh.Outgoing,msgCh.NumOnline)
	

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
