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
	



}