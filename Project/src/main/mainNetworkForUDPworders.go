package main

import (
	"assigner"
	def "definitions"
	"fmt"
	"fsm"
	"hardware"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
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
		FloorLamp:      make(chan int, 2),
		DoorLamp:       make(chan bool),
		BtnPressed:     make(chan def.ButtonPress, 10),
		DoorTimerReset: make(chan bool),
	}
	currentFloor := hardware.Init()
	time.Sleep(time.Millisecond * 500)
	go fsm.Init(eventCh, hardwareCh, currentFloor)
	go network.InitUDP(messageCh.Incoming, messageCh.Outgoing, eventCh.ElevatorPeerUpdate)
	go queue.RunBackup(messageCh.Outgoing)
	go EventHandler(eventCh, messageCh, hardwareCh)
	time.Sleep(time.Millisecond * 500)
	go assigner.CollectCosts(messageCh.CostReply)

	go safeKill()
	fmt.Println("Started main")
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
