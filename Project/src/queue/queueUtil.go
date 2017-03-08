package queue

import (
	def "definitions"
	//"log"
	//"fmt"
	"time"
)

type OrderInfo struct {
	Status bool
	Addr   string
	Timer  *time.Timer
}

type QueueType struct {
	Matrix [def.NumFloors][def.NumButtons]OrderInfo
}

var queue QueueType
var takeBackup = make(chan bool, 10)
var NewOrder = make(chan bool, 10)
var OrderTimeoutChan = make(chan def.ButtonPress, 10)
var LightUpdate = make(chan def.LightUpdate, 10)

func AddOrder(floor, btn int, addr string) {
	if queue.hasOrder(floor, btn) == false {
		queue.setOrder(floor, btn, OrderInfo{true, addr, nil})
		if addr == def.LocalElevatorId {
			NewOrder <- true
		} else {
			go queue.startTimer(floor, btn)
		}
	}
}

func RemoveOrder(floor, btn int, addr string) {
	if queue.Matrix[floor][btn].Addr == addr {
		queue.setOrder(floor, btn, OrderInfo{false, "", nil})
		queue.stopTimer(floor, btn)
	}
}

func OrderCompleted(floor, dir int, outgoingMsgCh chan<- def.Message) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if queue.Matrix[floor][btn].Addr == def.LocalElevatorId {
			if btn == def.BtnInside {
				RemoveOrder(floor, btn, def.LocalElevatorId)
			} else {
				outgoingMsgCh <- def.Message{Category: def.CompleteOrder, Floor: floor, Button: btn, Addr: def.LocalElevatorId}
			}
		}
	}
}

func ReassignOrder(floor, btn int, outgoingMsg chan<- def.Message, addr string) {
	RemoveOrder(floor, btn, addr)
	outgoingMsg <- def.Message{Category: def.NewOrder, Floor: floor, Button: btn}
}

func ReassignOrdersFromDeadElevator(addr string, outgoingMsgCh chan<- def.Message) {
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if queue.Matrix[floor][btn].Addr == addr {
				ReassignOrder(floor, btn, outgoingMsgCh, addr)
			}
		}
	}
}

func (q *QueueType) setOrder(floor, btn int, order OrderInfo) {
	q.Matrix[floor][btn] = order
	LightUpdate <- def.LightUpdate{Floor: floor, Button: btn, UpdateTo: order.Status}
	takeBackup <- true
	printQueue()
}

func (q *QueueType) startTimer(floor, btn int) {

	q.Matrix[floor][btn].Timer = time.NewTimer(def.ElevatorOrderTimeoutDuration)
	<-q.Matrix[floor][btn].Timer.C
	if q.Matrix[floor][btn].Status {
		OrderTimeoutChan <- def.ButtonPress{floor, btn}

	}
}

func (q *QueueType) stopTimer(floor, btn int) {
	if q.Matrix[floor][btn].Timer != nil {
		q.Matrix[floor][btn].Timer.Stop()
	}
}
