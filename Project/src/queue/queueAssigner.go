package queue

import (
	def "definitions"
	"fmt"
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
		fmt.Println("replylist:", replyList)
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
			fmt.Println("Best elevator is:\t", bestElevator)
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
