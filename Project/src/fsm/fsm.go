package fsm

import (
	def "definitions"
	"log"
	"queue"
	"time"
)

//Defining elevator behaviour
const (
	idle int = iota
	moving
	doorOpen
)

var Elevator struct {
	Floor     int
	Dir       int
	Behaviour int
}

func Init(eventCh def.EventChan, hardwareCh def.HardwareChan, startFloor int) {
	Elevator.Behaviour = idle
	Elevator.Dir = def.DirIdle
	Elevator.Floor = startFloor
	go doorTimer(hardwareCh.DoorTimerReset, eventCh.DoorTimeout)
	log.Println(def.ColG, "FSM initialized.", def.ColN)
}

func OnNewOrder(OutgoingMsg chan<- def.Message, hardwareCh def.HardwareChan) {
	switch Elevator.Behaviour {
	case doorOpen:
		if queue.ShouldStop(Elevator.Floor, Elevator.Dir) {
			hardwareCh.DoorTimerReset <- true
			queue.OrderCompleted(Elevator.Floor, Elevator.Dir, OutgoingMsg)
		}
	case moving:
		//Do nothing
	case idle:
		Elevator.Dir = queue.ChooseDirection(Elevator.Floor, Elevator.Dir)
		if Elevator.Dir == def.DirIdle {
			hardwareCh.DoorLamp <- true
			hardwareCh.DoorTimerReset <- true
			queue.OrderCompleted(Elevator.Floor, Elevator.Dir, OutgoingMsg)
			Elevator.Behaviour = doorOpen
		} else {
			hardwareCh.MotorDir <- Elevator.Dir
			Elevator.Behaviour = moving
		}
	}
}

func OnFloorArrival(hardwareCh def.HardwareChan, OutgoingMsg chan<- def.Message, newFloor int) {
	Elevator.Floor = newFloor
	hardwareCh.FloorLamp <- Elevator.Floor
	switch Elevator.Behaviour {
	case moving:
		if queue.ShouldStop(newFloor, Elevator.Dir) {
			tempDir := Elevator.Dir
			hardwareCh.MotorDir <- def.DirIdle
			hardwareCh.DoorLamp <- true
			hardwareCh.DoorTimerReset <- true
			queue.OrderCompleted(Elevator.Floor, tempDir, OutgoingMsg)
			Elevator.Behaviour = doorOpen
		}
	}
}

func OnDoorTimeout(hardwareCh def.HardwareChan) {
	switch Elevator.Behaviour {
	case doorOpen:
		Elevator.Dir = queue.ChooseDirection(Elevator.Floor, Elevator.Dir)
		hardwareCh.DoorLamp <- false
		hardwareCh.MotorDir <- Elevator.Dir
		if Elevator.Dir == def.DirIdle {
			Elevator.Behaviour = idle
		} else {
			Elevator.Behaviour = moving
		}
	}
}

//Private
func doorTimer(reset <-chan bool, timeout chan<- bool) {
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-reset:
			timer.Reset(def.ElevatorDoorTimeoutDuration)
		case <-timer.C:
			timer.Stop()
			timeout <- true
		}
	}
}
