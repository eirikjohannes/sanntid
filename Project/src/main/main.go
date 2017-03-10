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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)


func main() {

	/*queue kjører backup hver gang det kommer en ny ordre, OG hver gang en ordre fjernes fra køen.
	Dersom heisen blir stående uten ordre vil den time ut og restarte, dette burde ikke være noe problem.*/

	_, filepath1, _, _ := runtime.Caller(0)
	currentDirectory, _ := filepath.Split(filepath1)
	storageFile := currentDirectory + "/elevatorbackup.dat"
	
	fmt.Println("''''___Starting a new process____''''''")
	backup, err := os.Open(storageFile)
	if err != nil {
		tempFile, _ := os.Create(storageFile)
		tempFile.Close()
	}

	//An infinite loop that is only exited once the file has
	//not been written to for the last ElevatorResetTimeout. Once there has been no new orders for reset interval or no completion of orders.
	for true {
		fmt.Println("Looking for new entries in file...")
		fileStatus, err := os.Stat(storageFile)
		if err != nil {
			fmt.Println("Unrecoverable error...", err.Error())
			os.Exit(0)
		}
		if time.Now().After(fileStatus.ModTime().Add(def.ElevatorResetTimeout)) {
			break
		}
		time.Sleep(def.ElevatorOrderTimeoutDuration)
	}

	fmt.Println("____________Starts the main process_______")

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", filepath1)
	cmd.Start()
	log.Println(def.ColG, "****A new elevator spawns****",def.ColN)
	elevmain()
	
}

func elevmain() {
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

