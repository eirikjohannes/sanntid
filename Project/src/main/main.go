package main

import (
	"bufio"
	def "definitions"
	"fmt"
	"fsm"
	"hardware"
	"io/ioutil"
	"log"
	"network"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"queue2"
	"runtime"
	"strconv"
	"time"
)

func main() {

	/*queue kjører backup hver gang det kommer en ny ordre, OG hver gang en ordre fjernes fra køen.
	Dersom heisen blir stående uten ordre vil den time ut og restarte, dette burde ikke være noe problem.*/

	_, filepath1, _, _ := runtime.Caller(0)
	currentDirectory, _ := filepath.Split(filepath1)
	storageFile := currentDirectory + "aliveElev.dat"
	var counter = 0
	const aliveTimeout = 3 * time.Second
	const aliveInterval = 500 * time.Millisecond

	fmt.Println("''''___Starting a new process____''''''")
	backup, err := os.Open(storageFile)
	if err != nil {
		backup, _ := os.Create(storageFile)
		backup.Close()
	}

	//An infinite loop that is only exited once the file has
	//not been written to for the last ElevatorResetTimeout. Once there has been no new orders for reset interval or no completion of orders.
	for true {
		fmt.Println("Monitoring elevator")
		fileStatus, err := os.Stat(storageFile)
		if err != nil {
			fmt.Println("Unrecoverable error...", err.Error())
			os.Exit(0)
		}
		if time.Now().After(fileStatus.ModTime().Add(aliveTimeout)) {
			break
		}
		time.Sleep(aliveInterval)
	}

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", filepath1)
	cmd.Start()
	log.Println(def.ColG, "****A new elevator spawns****", def.ColN)

	backupReader := bufio.NewReader(backup)
	buffer, _ := backupReader.Peek(8)

	counter, _ = strconv.Atoi(string(buffer))
	backup.Close()
	backup, _ = os.Create(storageFile)

	defer func() {
		err := backup.Close()
		if err != nil {
			panic(err)
		}
	}()

	go elevmain()
	for true {
		counter++
		counterLineToBackup := strconv.Itoa(counter)
		ioutil.WriteFile(storageFile, []byte(counterLineToBackup), 0644)
		time.Sleep(aliveInterval)
	}
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
	go queue2.RunBackup(messageCh.Outgoing)
	go EventHandler(eventCh, messageCh, hardwareCh)
	go queue2.CollectCosts(messageCh.CostReply)

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

func EventHandler(eventCh def.EventChan, messageCh def.MessageChan, hardwareCh def.HardwareChan) {
	go eventButtonPressed(hardwareCh.BtnPressed)
	go eventElevatorAtFloor(eventCh.FloorReached)

	for {

		select {
		case btnPress := <-hardwareCh.BtnPressed:
			handleBtnPress(btnPress, messageCh.Outgoing)
		case incomingMsg := <-messageCh.Incoming:
			go sortAndHandleMessage(incomingMsg, messageCh)
		case btnLightUpdate := <-queue2.LightUpdate:
			log.Println(def.ColW, "Light update", def.ColN)
			hardware.SetBtnLamp(btnLightUpdate)
		case orderTimeout := <-queue2.OrderTimeoutChan:
			queue2.ReassignOrder(orderTimeout.Floor, orderTimeout.Button, messageCh.Outgoing)
		case motorDir := <-hardwareCh.MotorDir:
			hardware.SetMotorDir(motorDir)
		case floorLamp := <-hardwareCh.FloorLamp:
			hardware.SetFloorLamp(floorLamp)
		case doorLamp := <-hardwareCh.DoorLamp:
			hardware.SetDoorLamp(doorLamp)
		case <-queue2.NewOrder:
			log.Println(def.ColR, "Event: New order", def.ColN)
			fsm.OnNewOrder(messageCh.Outgoing, hardwareCh)
		case currFloor := <-eventCh.FloorReached:
			fsm.OnFloorArrival(hardwareCh, messageCh.Outgoing, currFloor)
		case <-eventCh.DoorTimeout:
			fsm.OnDoorTimeout(hardwareCh)
		case elevatorPeerUpdate := <-eventCh.ElevatorPeerUpdate:
			def.OnlineElevators = elevatorPeerUpdate.NumOnline
			if def.OnlineElevators == 0 {
				def.LocalElevatorId = "DISCONNECTED"
				log.Println(def.ColR, "Not connected to network, assigned DISCONNECTED", def.ColN)
			} else if def.LocalElevatorId == "DISCONNECTED" {
				network.AssignId()
				deadElevator := []string{"DISCONNECTED"}
				handleDeadElevator(deadElevator, messageCh.Outgoing)
			}
			if len(elevatorPeerUpdate.Lost) != 0 {
				handleDeadElevator(elevatorPeerUpdate.Lost, messageCh.Outgoing)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func eventButtonPressed(hardwareCh chan<- def.ButtonPress) {
	var buttonStateArray [def.NumFloors][def.NumButtons]bool

	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {

				if (floor == 0 && btn == def.BtnDown) || (floor == def.NumFloors-1 && btn == def.BtnUp) {

					continue
					//"Invalid operation", do nothing
				}
				if hardware.ReadButton(floor, btn) {

					if !(buttonStateArray[floor][btn]) {
						hardwareCh <- def.ButtonPress{Floor: floor, Button: btn}
					}
					buttonStateArray[floor][btn] = true
				} else {
					buttonStateArray[floor][btn] = false
				}
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func eventElevatorAtFloor(tempCh chan<- int) {
	var FloorReached = -7
	var prevFloor = -10
	for {
		if hardware.GetFloor() != -1 {
			FloorReached = hardware.GetFloor()
			if prevFloor != FloorReached {
				tempCh <- FloorReached
				prevFloor = FloorReached
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func sortAndHandleMessage(incomingMsg def.Message, messageCh def.MessageChan) {
	switch incomingMsg.Category {
	case def.NewOrder:
		log.Println(def.ColC, "New order incomming!", def.ColN)
		cost := queue2.CalculateCost(fsm.Elevator.Dir, hardware.GetFloor(), fsm.Elevator.Floor, incomingMsg.Floor, incomingMsg.Button)
		messageCh.Outgoing <- def.Message{Category: def.Cost, Floor: incomingMsg.Floor, Button: incomingMsg.Button, Cost: cost, Addr: def.LocalElevatorId}
	case def.CompleteOrder:
		queue2.RemoveOrder(incomingMsg.Floor, incomingMsg.Button)
		log.Println(def.ColG, "Order is completed", def.ColN)
	case def.Cost:
		log.Println(def.ColC, "Cost reply recieved as event", def.ColN)
		messageCh.CostReply <- incomingMsg
	}
}

func handleBtnPress(btnPress def.ButtonPress, outgoingMsg chan<- def.Message) {
	if btnPress.Button == def.BtnInside {
		queue2.AddOrder(btnPress.Floor, btnPress.Button, def.LocalElevatorId)
	} else {
		outgoingMsg <- def.Message{Category: def.NewOrder, Floor: btnPress.Floor, Button: btnPress.Button, Cost: 0, Addr: def.LocalElevatorId}
	}
}

func handleDeadElevator(address []string, outgoingMsg chan<- def.Message) {
	for i := 0; i < len(address); i++ {
		queue2.ReassignOrdersFromDeadElevator(address[i], outgoingMsg)
	}
}
