package queue

import (
	def "definitions"
	"log"
	"time"
)

type OrderInfo struct {

	Status  bool
	Addr    string
	Timer   *time.Timer
}

type QueueType struct {
	Matrix[def.Numfloors][def.NumButtons]OrderInfo
}

var queue QueueType
var takeBackup       = make(chan bool, 10)
var NewOrder         = make(chan bool, 10)
var OrderTimeoutChan = make(chan def.BtnPress, 10)
var LightUpdate      = make(chan def.LightUpdate, 10)

func AddOrder(floor, btn int, addr string) {
	if queue.hasOrder(floor,btn) == false {
		queue.setOrder(floor, btn, OrderInfo{true, addr, nil})
		if addr == def.LocalIP {
			NewOrder <- true
		}else{
			go queue.startTimer(floor,btn)
		}
	}
}

func RemoveOrder(floor, btn int) {
	queue.setOrder(floor, btn, OrderInfo{false, "", nil})
	queue.stopTimer(floor, btn)
}


func OrderCompleted(floor int, outgoingMsgCh chan<- def.message) {
	for btn := 0; btn < def.NumButtons; btn++ {
		if queue.Matrix[floor][btn].Addr == def.LocalIP {
			if btn == def.BtnInside {
				RemoveOrder(floor, btn)
			}else{
				outgoingMsgCh <- def.Message{def.CompleteOrder, floor, btn}
			}
		}
	}
}

func ReassignOrder(floor, btn int, outgoingMsg chan<- def.Message) {
	RemoveOrder(floor, btn)
	//log jepp jepp
	outgoingMsg <- def.Message{def.NewOrder, floor, btn}
}

func ReassignOrdersFromDeadElevator(addr string, outgoingMsgCh chan<- def.Message) {
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if queue.Matrix[floor][btn].Addr == addr {
				ReassignOrder(floor, btn, outgoingMsgCh)
			}
		}
	}
}

func (q *QueueType) setOrder(floor, btn int, order OrderInfo) {
	q.Matrix[floor][btn] = order
	LightUpdate <- def.LightUpdate{Floor: floor, Button: btn, UpdateTo: order.Status}
	takeBackup <- true
	//printQueue()
}

func (q *QueueType) startTimer(floor, btn int) {

	q.Matrix[floor][btn].Timer = time.NewTimer(def.ElevatorOrderTimeoutDuration)
	<-q.Matrix[floor][btn].Timer.C
	if q.Matrix[floor][btn].Status {
		OrderTimeoutChan <- def.BtnPress{floor, btn}

	}
}

func (q *QueueType) stopTimer(floor, btn int) {
	if q.Matrix[floor][btn].Timer != nil {
		q.Matrix[floor][btn].Timer.Stop()
	}
}
