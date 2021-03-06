package network

import (
	def "definitions"
	"network/bcast"
	"network/localip"
	"network/peers"
	"fmt"
	"log"
)

func InitUDP(incomingMsg chan def.Message, outgoingMsg chan def.Message, ElevatorPeerUpdateCh chan def.PeerUpdate) {

	def.OnlineElevators = 1
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
		def.OnlineElevators = 0
	}

	def.LocalElevatorId = localIP //fmt.Sprintf(localIP) //fmt.Sprintf("%s-%d", localIP, os.Getpid())
	//Initialize peerServer that handles alive and lost elevators

	peerTxEnable := make(chan bool)
	go peers.Transmitter(def.UDPPort, def.LocalElevatorId, peerTxEnable)
	go peers.Reciever(def.UDPPort, ElevatorPeerUpdateCh)

	//Initialize transmit and recieve servers for UDP messages
	msgRx := make(chan def.Message)
	msgTx := make(chan def.Message)
	go bcast.Reciever(def.UDPPort, msgRx) //WHyTHE FUCK wont it find .Reciever
	go bcast.Transmitter(def.UDPPort, msgTx)

	go forwardIncoming(incomingMsg, msgRx)
	go forwardOutgoing(outgoingMsg, msgRx, msgTx)

	log.Println(def.ColG, "Network initialized - IP: ", def.LocalElevatorId, def.ColN)
}

func AssignId() {
	tempID, err := localip.LocalIP()
	fmt.Println(err)
	log.Println(def.ColG, "New ID aqcuired: ", tempID, def.ColN)
	def.LocalElevatorId = tempID

}
func forwardIncoming(incomingMsg chan<- def.Message, msgRx <-chan def.Message) {
	for {
		msg := <-msgRx
		incomingMsg <- msg
	}
}

func forwardOutgoing(outgoingMsg <-chan def.Message, msgRx chan def.Message, msgTx chan<- def.Message) {
	for {
		msg := <-outgoingMsg
		if def.OnlineElevators == 0 {
			msgRx <- msg
		} else {
			msgTx <- msg
		}

	}
}
