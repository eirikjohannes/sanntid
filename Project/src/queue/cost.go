package queue

import def "definitions"

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
			totCost += 10
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
