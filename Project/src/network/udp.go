
package network

import (
	def "definitions"
	"network/bcast"
	"network/localip"
	"network/peers"
	//"flag"
	"fmt"
	"os"
	//"time"
	"log"

)




func InitUDP(incomingMsg chan def.Message, outgoingMsg chan def.Message, ElevatorPeerUpdateCh chan def.PeerUpdate) {
	
	//var elevatorId string
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	def.LocalElevatorId = localIP.String()//fmt.Sprintf(localIP) //fmt.Sprintf("%s-%d", localIP, os.Getpid())

	//Initialize peerServer that handles alive and lost elevators
	

	peerTxEnable := make(chan bool)
	go peers.Transmitter(def.UDPPort, def.LocalElevatorId, peerTxEnable)
	go peers.Reciever(def.UDPPort, ElevatorPeerUpdateCh, NumOnlineCh)

	//Initialize transmit and recieve servers for UDP messages
	msgRx := make(chan def.Message)
	msgTx := make(chan def.Message)
	go bcast.Reciever(def.UDPPort, msgRx)
	go bcast.Transmitter(def.UDPPort, msgTx)
	
	
	go forwardIncoming(incomingMsg,msgRx)
	go forwardOutgoing(outgoingMsg,msgTx)

	log.Println(def.ColG, "Network initialized - IP: ", def.LocalIP, def.ColN)
}


func forwardIncoming(incomingMsg chan<- def.Message, msgRx <-chan def.Message){
	for{
		msg:=<-msgRx
		incomingMsg<-msg
	}
}

func forwardOutgoing(outgoingMsg <-chan def.Message, msgTx chan<- def.Message){
	for{
		msg:=<-outgoingMsg
		msgTx <- msg
	}	
}


