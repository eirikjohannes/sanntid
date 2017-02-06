package udp

import (
	"./bcast"
	"./localip"
	"./peers"
	def "definitions"
	"flag"
	"fmt"
	"os"
	"time"
)



func InitUDP()
{

	var elevatorId string
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	elevatorId = fmt.Sprintf("%s-%d", localIP, os.Getpid())

	elevatorPeerUpdateChannel := make(chan peers.PeerUpdate)

	go peers.Transmitter(def.UDPPort, elevatorId, True)
	go peers.Reciever(def.UDPPort, elevatorPeerUpdateChannel)



	orderTx := make(chan def.ElevatorOrder)
	orderRx := make(chan def.ElevatorOrder)
	aliveTx := make(chan def.ElevatorAliveMessage)
	aliveRx := make(chan def.ElevatorAliveMessage)

	go bcast.Transmitter(def.UDPPort, orderTx, aliveTx)
	go bcast.Reciever(def. UDPPort, orderRx, aliveRx)
	
	//Decide which queues should transport the orders. Probably a channel that will communicate with the other threads. This should communicate with queue only....?

	runUDP()
}


Channel to create in main: 
NetworkToQueueOrderChannel := make(chan def.ElevatorOrder)
NetworkToQueueAliveChannel := make(chan def.ElevatorAliveMessage)
QueueToNetworkOrderChannel := make(chan def.ElevatorOrder)
QueueToNetworkAliveChannel := make(chan def.ElevatorAliveMessage)

//Add order to queue tar imot fra nettverket og formidler til queue. Når queue bekrefter, returneres en positiv verdi som fungerer som en ack. Deretter sender nettverksmodulen ut ACK.
func addOrderToQueue(ElevatorOrder order) int {
	order := <- orderRx
	networkToQueueOrderChannel<-order
	addedOrder:=<-QueueToNetworkOrderChannel 
	for addedOrder.Ack!=True{ //or orderID order.)
		addedOrder :=<-QueueToNetworkOrderChannel
	}
	//We now have sucessfully added the order to the queue. Send ack onto network.
	orderTx<-addedOrder
	//Hvis det kommer en ny ting i orderRx, send til kanalen som kommuniserer mellom kø og netwrok. vent på ack fra queue og send deretter ACK på nettverket.
}

func distributeOrderToNetwork(ElevatorOrder order){
	//Hvis ny ordre ligger i kanalen, ta imot og legg i orderTx kanalen. 

	//vent på bekreftelse fra alle kjente peers.
}

func runUDP(){

}