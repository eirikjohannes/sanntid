package network

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

var LocalElevatorId string
var MapOfAckedOrders map[def.ElevatorOrder]int
var ListOfPeers[]string

func InitUDP()
{
	MapOfAckedOrders = make(map[def.ElevatorOrder]int)

	//var elevatorId string
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	LocalElevatorId = fmt.Sprintf("%s-%d", localIP, os.Getpid())

	ElevatorPeerUpdateChannel := make(chan peers.PeerUpdate)

	go peers.Transmitter(def.UDPPort, LocalElevatorId, True)
	go peers.Reciever(def.UDPPort, ElevatorPeerUpdateChannel)



	orderTx := make(chan def.ElevatorOrder)
	orderRx := make(chan def.ElevatorOrder)
	aliveTx := make(chan def.ElevatorAliveMessage)
	aliveRx := make(chan def.ElevatorAliveMessage)

	go bcast.Transmitter(def.UDPPort, orderTx, aliveTx)
	go bcast.Reciever(def. UDPPort, orderRx, aliveRx)
	
	//Decide which queues should transport the orders. Probably a channel that will communicate with the other threads. This should communicate with queue only....?
	go addOrderToQueue()
	go distributeOrderToNetwork()
	
}

func GetElevatorId() string{
	return LocalElevatorId;
}
/*Channel to create in main: 
NetworkToQueueOrderChannel := make(chan def.ElevatorOrder,10)
NetworkToQueueAliveChannel := make(chan def.ElevatorAliveMessage,10)
QueueToNetworkOrderChannel := make(chan def.ElevatorOrder,10)
QueueToNetworkAliveChannel := make(chan def.ElevatorAliveMessage,10)

*/

//Add order to queue tar imot fra nettverket og formidler til queue. Når queue bekrefter, returneres en positiv verdi som fungerer som en ack. Deretter sender nettverksmodulen ut ACK.
func addOrderToQueue(){
	//def.ElevatorOrder order
	for{
		order := <- orderRx
		select{
			case !order.Ack			
				NetworkToQueueOrderChannel<-order
				addedOrder:=<-QueueToNetworkOrderChannel 
				for addedOrder.Ack!=True{ //or orderID order.)
					addedOrder :=<-QueueToNetworkOrderChannel
				}
				//We now have sucessfully added the order to the queue. Send ack onto network.
				orderTx<-addedOrder
			case order.Ack
				//Order is acked from another source
				//If the order was sent from this elevator, append ACK to list.
				if (order.ElevatorID == LocalElevatorId){
					orderAlreadyInMap:=addAck(&order)
					// if (orderAlreadyInMap){
					// 	orderTx<-order //resend to network
					// 	order.Ack=0
					// 	NetworkToQueueOrderChannel<-order //resends order to queue to add order again.
					// }
				}
		}//Hvis det kommer en ny ting i orderRx, send til kanalen som kommuniserer mellom kø og netwrok. vent på ack fra queue og send deretter ACK på nettverket.
	}
}

func distributeOrderToNetwork(ElevatorOrder order){
	//Hvis ny ordre ligger i kanalen, ta imot og legg i orderTx kanalen. 
	for{
		
		orderToDistribute := <-QueueToNetworkOrderChannel
		switch i:=addAck(orderToDistribute); i {
			case 1:
				fmt.Println("Order was distributed, but already in mapOfAckedOrders. ElevatorID: "+orderToDistribute.ElevatorId)
			case 0:
				fmt.Println("Order was distributed and did not exist in mapOfAckedOrders. ElevatorID: "+orderToDistribute.ElevatorId)
		}
		orderTx<-orderToDistribute
	}
	//vent på bekreftelse fra alle kjente peers.
}

func addAck(orderToCheck *def.ElevatorOrder) int{
	var orderExists bool=false
	var tempValue int =0
	tempValue, orderExists:=MapOfAckedOrders[*orderToCheck]


	//printmap
	fmt.Println("Before if(OrderExists)\n")
	for key, value := range MapOfAckedOrders{
		fmt.Println("Key:",key, "NumberofAcks: %i\n", value)
	}

	if (orderExists){
		if (tempValue+1==len(peers.GetNumberOfPeers)){
			delete (MapOfAckedOrders, *orderToCheck)
		}
		else{
			MapOfAckedOrders[*orderToCheck]=(tempValue+1)
		}
		return 1
	}
	else{
		MapOfAckedOrders[*orderToCheck]=0;
		return 0
	}
	fmt.Println("After if(OrderExists)\n")
	for key, value := range MapOfAckedOrders{
		fmt.Println("Key:",key, "NumberofAcks: %i\n", value)
	}
	

}	