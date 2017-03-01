package main

import (
	def "definitions"
	"fmt"
	"hardware"
	"network"
	//"os"
	"time"
)

func main() {
	messageCh := def.MessageChan{
		Outgoing:  make(chan def.Message, 10),
		Incoming:  make(chan def.Message, 10),
		CostReply: make(chan def.Message, 10),
	}
	eventCh := def.EventChan{
		FloorReached:       make(chan int, 10),
		DoorTimeout:        make(chan bool, 10),
		ElevatorPeerUpdate: make(chan def.PeerUpdate, 2),
	}
	hardwareCh := def.HardwareChan{
		MotorDir:       make(chan int),
		FloorLamp:      make(chan int),
		DoorLamp:       make(chan bool),
		BtnPressed:     make(chan def.ButtonPress, 10),
		DoorTimerReset: make(chan bool),
	}
	go network.InitUDP(messageCh.Incoming, messageCh.Outgoing, eventCh.ElevatorPeerUpdate)
	go PrintOrder(messageCh.Incoming, messageCh.Outgoing)
	go hardware.Init()
	go EventHandler(eventCh, messageCh, hardwareCh)

	time.Sleep(time.Millisecond * 500)
	order := def.Message{Category: def.NewOrder, Floor: 1, Button: 2, Cost: 0, Addr: def.LocalElevatorId}

	orderDistributeTimer := time.NewTicker(time.Second * 2)
	for {
		<-orderDistributeTimer.C
		messageCh.Outgoing <- order
		messageCh.Outgoing <- def.Message{def.CompleteOrder, 1, 2, 0, def.LocalElevatorId}
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
			fmt.Println("The new order has ID %s", orderToPrint.Addr)
			fmt.Println("Floor; %i \t Button; %i \t", orderToPrint.Floor, orderToPrint.Button)
		case def.CompleteOrder:
			fmt.Println("\n\tCOMPLETED order!\n")
			fmt.Println("The completed order has ID %s", orderToPrint.Addr)
			fmt.Println("Floor; %i \t Button; %i \t", orderToPrint.Floor, orderToPrint.Button)
		case def.Cost:
			fmt.Println("Cost msg recieved")
		}

	}

}
