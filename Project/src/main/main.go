package main

import (
	def "definitions"
	"fsm"
	"hardware"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
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
		MotorDir:       make(chan int, 2),
		FloorLamp:      make(chan int, 2),
		DoorLamp:       make(chan bool, 2),
		BtnPressed:     make(chan def.ButtonPress, 10),
		DoorTimerReset: make(chan bool, 2),
	}
	currentFloor := hardware.Init()
	go fsm.Init(eventCh, hardwareCh, currentFloor)
	go network.InitUDP(messageCh.Incoming, messageCh.Outgoing, eventCh.ElevatorPeerUpdate)
	go queue.RunBackup(messageCh.Outgoing)
	go EventHandler(eventCh, messageCh, hardwareCh)
	go queue.CollectCosts(messageCh.CostReply)

	go safeKill()

	hold := make(chan bool)
	<-hold

}

func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hardware.SetMotorDir(def.DirIdle)
	log.Fatal(def.Col0, "User terminated program.", def.ColN)
}
