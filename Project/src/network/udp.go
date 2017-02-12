package network

import (
	"network/bcast"
	"network/localip"
	"network/peers"
	def "definitions"
	//"flag"
	"fmt"
	"os"
	//"time"
)

var LocalElevatorId string
var MapOfAckedOrders map[def.ElevatorOrder]int
var ListOfPeers[]string


func InitUDP(NetworkToQueueOrderChannel chan def.ElevatorOrder, QueueToNetworkOrderChannel chan def.ElevatorOrder, NetworkToQueueAliveChannel chan def.ElevatorAliveMessage, QueueToNetworkAliveChannel chan def.ElevatorAliveMessage){
	MapOfAckedOrders = make(map[def.ElevatorOrder]int)
	//var elevatorId string
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	LocalElevatorId = fmt.Sprintf("%s-%d", localIP, os.Getpid())

	ElevatorPeerUpdateChannel := make(chan peers.PeerUpdate)

	peerTxEnable :=make(chan bool)
	go peers.Transmitter(def.UDPPort, LocalElevatorId, peerTxEnable);
	go peers.Cheesedoodles(def.UDPPort, ElevatorPeerUpdateChannel);

	orderTx := make(chan def.ElevatorOrder)
	orderRx := make(chan def.ElevatorOrder)
	aliveTx := make(chan def.ElevatorAliveMessage)
	aliveRx := make(chan def.ElevatorAliveMessage)
	
	go bcast.Transmitter2(def.UDPPort, orderRx) 
	go bcast.Transmitter2(def.UDPPort, aliveRx)
	go bcast.Transmitter(def.UDPPort, orderTx)
	go bcast.Transmitter(def.UDPPort, aliveTx)

	
	//Decide which queues should transport the orders. Probably a channel that will communicate with the other threads. This should communicate with queue only....?
	go addOrderToQueue(orderTx,orderRx, QueueToNetworkOrderChannel, NetworkToQueueOrderChannel)
	go distributeOrderToNetwork(orderTx, QueueToNetworkOrderChannel)
	
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

//Add order to queue tar imot fra nettverket og formidler til queue. 
//Når queue bekrefter, returneres en positiv verdi som fungerer som en ack. Deretter sender nettverksmodulen ut ACK.
func addOrderToQueue(orderTx chan def.ElevatorOrder, orderRx chan def.ElevatorOrder, QueueToNetworkOrderChannel chan def.ElevatorOrder, NetworkToQueueOrderChannel chan def.ElevatorOrder){
	//def.ElevatorOrder order
	for{
		order := <- orderRx
		if(!order.Ack){
			fmt.Println("\nFrom addordertoquque: \n The order has not been acked. New order, send to queue.\n")			
			NetworkToQueueOrderChannel<-order
			addedOrder:=<-QueueToNetworkOrderChannel 
			fmt.Println("Order returned to network\n")
			for (addedOrder.Ack!=true){ //or orderID order.)
				addedOrder = <-QueueToNetworkOrderChannel
			}
			fmt.Println("Order is acked, transmitting w/ack to ontonetwork\n")
			//Possible deadlock of addORderToQUeue
			//We now have sucessfully added the order to the queue. Send order with ack onto network.
			orderTx<-addedOrder
		}else if (order.Ack){
			//Order is acked from another source
			//If the order was sent from this elevator, append ACK to list.
			if (order.ElevatorId == LocalElevatorId){
				addAck(&order)
				fmt.Println("\n from addOrderToQueue:\n order.ACk=1, order is originally sent to this computer.\n\n")
				// if (orderAlreadyInMap){
				// 	orderTx<-order //resend to network
				// 	order.Ack=0
				// 	NetworkToQueueOrderChannel<-order //resends order to queue to add order again.
				// }
			}
		}
	}//Hvis det kommer en ny ting i orderRx, send til kanalen som kommuniserer mellom kø og netwrok. vent på ack fra queue og send deretter ACK på nettverket.
}

func distributeOrderToNetwork(orderTx chan def.ElevatorOrder,QueueToNetworkOrderChannel chan def.ElevatorOrder){
	//Hvis ny ordre ligger i kanalen, ta imot og legg i orderTx kanalen. 
	for{
		
		orderToDistribute := <-QueueToNetworkOrderChannel
		switch i:=addAck(&orderToDistribute); i {
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
	tempValue, orderExists = MapOfAckedOrders[*orderToCheck]

	//printmap
	fmt.Println("Before if(OrderExists)\n")
	for key, value := range MapOfAckedOrders{
		fmt.Println("Key: %d \t NumberofAcks: %d \n",key, value)
	}

	if (orderExists){
		fmt.Println("\nNumber of peers: %d", peers.GetNumberOfPeers())
		if ((tempValue+1)==peers.GetNumberOfPeers()){
			delete (MapOfAckedOrders, *orderToCheck)
		}else{
			MapOfAckedOrders[*orderToCheck]=(tempValue+1)
		}
		return 1
	}else{
		MapOfAckedOrders[*orderToCheck]=0;
		return 0
	}
	fmt.Println("After if(OrderExists)\n")
	for key, value := range MapOfAckedOrders{
		fmt.Println("Key:",key, "NumberofAcks: %d\n", value)
	}
	return 0
}	


