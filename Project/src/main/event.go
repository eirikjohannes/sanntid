//Event.go fil for

package main

import (
	//"assigner"
	def "definitions"
	"fsm"
	"log"
	"queue"
)

var elevatorPeerUpdate def.PeerUpdate

func EventHandler(eventCh def.EventChan, messageCh def.MessageChan, hardwareCh def.HardwareChan) {

	//Initialiser go funksjonen for
	//Initialiser go ufnksjonen for
	//resten skal kunne initialiseres andre steder.

	go eventButtonPressed(hardwareCh.BtnPressed)
	go eventElevatorAtFloor(eventCh.FloorReached)

	for {
		select {
		case btnPress := <-hardwareCh.BtnPressed:
			handleBtnPress(btnPress, messageCh.Outgoing)
		case incomingMsg := <-messageCh.Incoming:
			go sortAndHandleMessage(incomingMsg, messageCh)
		case btnLightUpdate := <-queue.LightUpdate:
			log.Println(def.ColW, "Light update", def.ColN)
			hardware.SetBtnLamp(btnLightUpdate)
		case orderTimeout := <-queue.OrderTimeoutChan:
			queue.ReassignOrder(orderTimeout.Floor, orderTimeout.Button, messageCh.Outgoing)
		case motorDir := <-hardwareCh.MotorDir:
			hardware.SetMotorDir(motorDir)
		case floorLamp := <-hardwareCh.FloorLamp:
			hardware.SetFloorLamp(floorLamp)
		case doorLamp := <-hardwareCh.DoorLamp:
			hardware.SetDoorLamp(doorLamp)
		case <-queue.NewOrder:
			log.Println(def.ColW, "Event: New order", def.ColN)
			fsm.OnNewOrder(messageCh.Outgoing, hardwareCh)
		case currFloor := <-eventCh.FloorReached:
			fsm.OnFloorArrival(hardwareCh, messageCh.Outgoing, currFloor)
		case <-eventCh.DoorTimeout:
			fsm.OnDoorTimeout(hardwareCh)
		case elevatorPeerUpdate <- eventCh.ElevatorPeerUpdate:
			if elevatorPeerUpdate.Lost {
				handleDeadElevator(elevatorPeerUpdate.Lost, outgoingMsg, messageCh.NumOnline)
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func eventButtonPressed(hwCh chan<- def.BtnPressed) {
	var buttonStateArray [def.NumFloors][def.NumButtons]bool

	for {
		for floor := 0; floor < def.NumFloors; floor++ {
			for btn := 0; btn < def.NumButtons; btn++ {
				if (floor == 0 && btn == def.BtnHallDown) || (floor == def.NumFloors-1 && btn == def.BtnHallUp) {
					continue
					//"Invalid operation", do nothing
				}
				if hardware.ReadButton(floor, btn) {
					if !(buttonStateArray[floor][btn]) {
						hwCh <- def.ButtonPress{Floor: floor, Button: btn}
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
		messageCh.Outgoing <- def.Message{Category: def.Cost, Floor: incomingMsg.Floor, Button: incomingMsg.Button, Cost: cost}
	case def.CompleteOrder:
		queue.RemoveRequest(incomingMsg.Floor, incomingMsg.Button)
		log.Println(def.ColG, "Order is completed", def.ColN)
	case def.Cost:
		log.Println(def.ColC, "Cost reply recieved as event", ColN)
		messageCh.CostReply <- incomingMsg
	}
}

func handleBtnPress(btnPress def.ButtonPress, outgoingMsg chan<- def.Message) {
	if btnPress.Button == def.BtnCab {
		queue.AddOrder(btnPress.Floor, btnPress.Button, def.LocalElevatorId)
	} else {
		outgoingMsg <- def.Message{Category: def.NewRequest, Floor: btnPress.Floor, Button: btnPress.Button, Cost: 0, Addr: def.LocalElevatorId}
	}
}

//This is not finished, sanity check this
func handleDeadElevator(address []string, outgoingMsg chan<- def.Message) {
	for i = 0; i < len(address); i++ {
		queue.ReassignAllRequestsFrom(address[i], outgoingMsg)
	}
}
