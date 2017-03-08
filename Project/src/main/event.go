//Event.go fil for

package main

import (
	//"assigner"
	def "definitions"
	"fmt"
	"fsm"
	"hardware"
	"log"
	"queue"
	"time"
)

func EventHandler(eventCh def.EventChan, messageCh def.MessageChan, hardwareCh def.HardwareChan) {

	//Initialiser go funksjonen for
	//Initialiser go ufnksjonen for
	//resten skal kunne initialiseres andre steder.

	go eventButtonPressed(hardwareCh.BtnPressed)
	go eventElevatorAtFloor(eventCh.FloorReached)

	for {

		select {
		case btnPress := <-hardwareCh.BtnPressed:
			fmt.Println("Button is pressed")
			handleBtnPress(btnPress, messageCh.Outgoing)
		case incomingMsg := <-messageCh.Incoming:
			go sortAndHandleMessage(incomingMsg, messageCh)
		case btnLightUpdate := <-queue.LightUpdate:
			log.Println(def.ColW, "Light update", def.ColN)
			hardware.SetBtnLamp(btnLightUpdate)
		case orderTimeout := <-queue.OrderTimeoutChan:
			queue.ReassignOrder(orderTimeout.Floor, orderTimeout.Button, messageCh.Outgoing)
		case motorDir := <-hardwareCh.MotorDir:
			fmt.Println("Did not Got here")
			hardware.SetMotorDir(motorDir)
		case floorLamp := <-hardwareCh.FloorLamp:
			hardware.SetFloorLamp(floorLamp)
		case doorLamp := <-hardwareCh.DoorLamp:
			hardware.SetDoorLamp(doorLamp)
		case <-queue.NewOrder:
			fmt.Println(def.ColW, "Event: New order", def.ColN)
			fsm.OnNewOrder(messageCh.Outgoing, hardwareCh)
		case currFloor := <-eventCh.FloorReached:
			fmt.Println("got here and...")
			fsm.OnFloorArrival(hardwareCh, messageCh.Outgoing, currFloor)
			fmt.Println("hung")
		case <-eventCh.DoorTimeout:
			fsm.OnDoorTimeout(hardwareCh)
		case elevatorPeerUpdate := <-eventCh.ElevatorPeerUpdate:
			fmt.Println("elevatorPeerUpdate")
			def.OnlineElevators = elevatorPeerUpdate.NumOnline
			if len(elevatorPeerUpdate.Lost) != 0 {
				handleDeadElevator(elevatorPeerUpdate.Lost, messageCh.Outgoing)
			}
		}
		fmt.Println("Completed one select")
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
						fmt.Println("Looping")
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
	case def.Alive:
		//Kan droppes pga peer stufF??
		/*address:=incomingMsg.Addr
		if connection, exist := onlineElevatorMap[address]; exist{
			connection.Timer.Reset(def.ElevTimeoutDuration)
			//Check how the connection type something works
		} else{
			newConnection := def.UdpConnection{Addr: addr, Timer: time.NewTimer(def.ElevTimeoutDuration)}
			onlineElevatorMap[addr] = newConnection
			msgCh.NumOnline <- len(onlineElevatorMap)
			go connectionTimer(&newConnection, msgCh.Outgoing, msgCh.NumOnline)
			log.Println(def.ColG, "New elevator: ", addr, " | Number online: ", len(onlineElevatorMap), def.ColN)
		}*/
	case def.NewOrder:
		log.Println(def.ColC, "New order incomming!", def.ColN)
		cost := queue.CalculateCost(fsm.Elevator.Dir, hardware.GetFloor(), fsm.Elevator.Floor, incomingMsg.Floor, incomingMsg.Button)
		messageCh.Outgoing <- def.Message{Category: def.Cost, Floor: incomingMsg.Floor, Button: incomingMsg.Button, Cost: cost, Addr: def.LocalElevatorId}
	case def.CompleteOrder:
		queue.RemoveOrder(incomingMsg.Floor, incomingMsg.Button)
		log.Println(def.ColG, "Order is completed", def.ColN)
	case def.Cost:
		log.Println(def.ColC, "Cost reply recieved as event", def.ColN)
		messageCh.CostReply <- incomingMsg
	}
}

func handleBtnPress(btnPress def.ButtonPress, outgoingMsg chan<- def.Message) {
	if btnPress.Button == def.BtnInside {
		queue.AddOrder(btnPress.Floor, btnPress.Button, def.LocalElevatorId)
	} else {
		outgoingMsg <- def.Message{Category: def.NewOrder, Floor: btnPress.Floor, Button: btnPress.Button, Cost: 0, Addr: def.LocalElevatorId}
	}
}

//This is not finished, sanity check this
func handleDeadElevator(address []string, outgoingMsg chan<- def.Message) {
	for i := 0; i < len(address); i++ {
		queue.ReassignOrdersFromDeadElevator(address[i], outgoingMsg)
	}
}
