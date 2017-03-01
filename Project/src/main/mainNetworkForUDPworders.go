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
		Outgoing:  make(chan def.Message, 10),
		Incoming:  make(chan def.Message, 10),
		CostReply: make(chan def.Message, 10),
	}
	eventCh := def.EventChan{
		ElevatorPeerUpdate: make(chan def.PeerUpdate, 2),
	}
	go network.InitUDP(msgCh.Incoming, msgCh.Outgoing, eventCh.ElevatorPeerUpdate)
	go PrintOrder(msgCh.Incoming, msgCh.Outgoing)

	order := def.Message{Category: def.NewOrder, Floor: 2, Button: 2, Cost: 0, Addr: def.LocalElevatorId}

	orderDistributeTimer := time.NewTicker(time.Second * 2)
	for {
		<-orderDistributeTimer.C
		msgCh.Outgoing <- order
	}

}

func PrintOrder(incomingMsg chan def.Message, outgoingMsg chan def.Message) {
	for {
		orderToPrint := <-incomingMsg
		switch orderToPrint.Category {
		case def.Alive:
			fmt.Println("Alive message recieved")
		case def.NewOrder:
			fmt.Println("NEW order recieved! Fantastic news\n")
			fmt.Println("The new order has ID %s \n", orderToPrint.Addr)
			fmt.Println("Floor; %i \t Button; %i \t", orderToPrint.Floor, orderToPrint.Button)
		case def.CompleteOrder:
			fmt.Println("\n\tCOMPLETED order!\n")
			fmt.Println("The completed order has ID %s \n", orderToPrint.Addr)
			fmt.Println("Floor; %i \t Button; %i \t", orderToPrint.Floor, orderToPrint.Button)
		case def.Cost:
			fmt.Println("Cost msg recieved")
		}

	}

}
