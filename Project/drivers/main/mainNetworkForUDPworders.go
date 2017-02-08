package main

import (
	def "definitions"
	"fmt"
	"network"
	"os"
	"time"
)

func main() {
	tick := time.Tick(1 * time.Second)
	initChannels()
	network.InitUDP()

	newOrder:=def.ElevatorOrder
	newOrder.Floor=3
	newOrder.Btn=1
	newOrder.ElevatorId=network.GetElevatorId()
	newOrder.Ack=0
	go printOrders()
	for{
		select{
			case <- tick
				QueueToNetworkOrderChannel <-newOrder
				newOrder.Floor--			
		}
	
	}
	
}

func initChannels() {
	NetworkToQueueOrderChannel := make(chan def.ElevatorOrder, 10)
	NetworkToQueueAliveChannel := make(chan def.ElevatorAliveMessage, 10)
	QueueToNetworkOrderChannel := make(chan def.ElevatorOrder, 10)
	QueueToNetworkAliveChannel := make(chan def.ElevatorAliveMessage, 10)

}

func printOrder(){
	orderToPrint:=<-NetworkToQueueOrderChannel
	fmt.Println("NEW order recieved! Fantastic news\n")
	fmt.Println("The new order has ID %s \n",orderToPrint.ElevatorId)
	fmt.Println("Floor; %i \t Button; %i \t Ack: %t \n_________\n", orderToPrint.Floor, orderToPrint.Btn, orderToPrint.Ack )
}
