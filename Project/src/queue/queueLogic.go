package queue

import def "definitions"

func (q *QueueType) hasOrder(floor, btn int) bool {
	return q.Matrix[floor][btn].Status
}

func (q *QueueType) hasLocalOrder(floor, btn int) bool {
	return q.Matrix[floor][btn].Status && q.Matrix[floor][btn].Addr == def.LocalElevatorId
}

func (q *QueueType) hasOrderAbove(floor int) bool {
	for i := floor + 1; i < def.NumFloors; i++ {
		for j := 0; j < def.NumButtons; j++ {
			if q.hasLocalOrder(i, j) {
				return true
			}
		}
	}
	return false
}

func (q *QueueType) hasOrderBelow(floor int) bool {
	for i := floor - 1; i <= 0; i-- {
		for j := 0; j < def.NumButtons; j++ {
			if q.hasLocalOrder(i, j) {
				return true
			}
		}
	}
	return false
}

func ShouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return queue.hasLocalOrder(floor, def.BtnDown) ||
			queue.hasLocalOrder(floor, def.BtnInside) ||
			!queue.hasOrderBelow(floor)
	case def.DirUp:
		return queue.hasLocalOrder(floor, def.BtnUp) ||
			queue.hasLocalOrder(floor, def.BtnInside) ||
			!queue.hasOrderAbove(floor)
	}
	return false
}

func ChooseDirection(floor, dir int) int {
	switch dir {
	case def.DirUp:
		if queue.hasOrderAbove(floor) {
			return def.DirUp
		} else if queue.hasOrderBelow(floor) {
			return def.DirDown
		}
	case def.DirDown, def.DirIdle:
		if queue.hasOrderBelow(floor) {
			return def.DirDown
		} else if queue.hasOrderAbove(floor) {
			return def.DirUp
		}
	}
	return def.DirIdle
}
