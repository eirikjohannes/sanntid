package queue

import (
	def "definitions"
	"log"
	"time"
)

type reply struct {
	cost     int
	elevator string
}
type order struct {
	floor  int
	button int
	timer  *time.Timer
}

func CollectCosts(costReply <-chan def.Message) {
	orderMap := make(map[order][]reply)
	var timeout = make(chan *order)
	for {
		select {
		case message := <-costReply:
			handleCostReply(orderMap, message, def.OnlineElevators, timeout)
		case <-timeout:
			log.Println(def.ColR, "Not all costs received in time!", def.ColN)
			chooseBestElevator(orderMap, def.OnlineElevators, true)
		}
	}
}

func CalculateCost(currentDir, currentFloor, prevFloor, targetFloor, targetButton int) int{
	totalCost := 0
	dir       := currentDir
	targetDir := targetFloor-prevFloor


	if currentFloor == -1 {
		totalCost++
	} else if dir != def.DirIdle {
		totalCost += 2
	}
	if dir != def.DirIdle {
		if targetDir != dir {
			totalCost += 10
		}
	}

	if targetDir > 0 && dir == def.DirUp || dir == def.DirIdle {
		for floor := prevFloor; floor < targetFloor || floor == def.NumFloors; floor++ {
			if queue.hasLocalOrder(floor, targetButton) || queue.hasLocalOrder(floor, def.BtnInside) {
				totalCost++
			}
			totalCost++
		}
	}
	if targetDir < 0 && dir == def.DirDown || dir == def.DirIdle {
		for floor := prevFloor; floor > targetFloor || floor == 0; floor-- {
			if queue.hasLocalOrder(floor, targetButton) || queue.hasLocalOrder(floor, def.BtnInside) {
				totalCost++
			}
			totalCost++
		}
	}
	return totalCost
}

func handleCostReply(orderMap map[order][]reply, message def.Message, numOnline int, timeout chan *order) {
	newOrder := order{floor: message.Floor, button: message.Button}
	newReply := reply{cost: message.Cost, elevator: message.Addr}

	for existingOrder := range orderMap {
		if equal(existingOrder, newOrder) {
			newOrder = existingOrder
		}
	}
	if replyList, exist := orderMap[newOrder]; exist {
		found := false
		for _, reply := range replyList {
			if reply == newReply {
				found = true
			}
		}
		if found == false {
			orderMap[newOrder] = append(orderMap[newOrder], newReply)
			newOrder.timer.Reset(def.CostReplyTimeoutDuration)
		}
	} else {
		newOrder.timer = time.NewTimer(def.CostReplyTimeoutDuration)
		orderMap[newOrder] = []reply{newReply}
		go costTimer(&newOrder, timeout)
	}
	chooseBestElevator(orderMap, numOnline, false)
}

func chooseBestElevator(orderMap map[order][]reply, numOnline int, isTimeout bool) {
	var bestElevator string

	for order, replyList := range orderMap {
		if numOnline == 0 {
			bestElevator = def.LocalElevatorId
			AddOrder(order.floor, order.button, bestElevator)
			order.timer.Stop()
			delete(orderMap, order)
		} else if len(replyList) == numOnline || isTimeout {
			lowestCost := 9001
			for _, reply := range replyList {
				if reply.cost < lowestCost {
					lowestCost = reply.cost
					bestElevator = reply.elevator
				} else if reply.cost == lowestCost {
					//choose elevator with lowest IP on equal cost
					if reply.elevator < bestElevator {
						bestElevator = reply.elevator
					}
				}
			}
			AddOrder(order.floor, order.button, bestElevator)
			order.timer.Stop()
			delete(orderMap, order)
		}
	}
}

func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
}

func costTimer(newOrder *order, timeout chan<- *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}
