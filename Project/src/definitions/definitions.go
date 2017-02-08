/*Definition file for elevator project
Group 13, Einar Henriksen, Eirik Larsen
*/

package definitions

import "time"

const NumFloors = 4
const NumButtons = 3
const ElevatorDoorTimeoutDuration = 2 * time.Second
const ElevatorTimeoutDuration = 5 * time.Second
const AliveMessageInterval = 500 * time.Millisecond
const UDPPort = 13131


const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
    BtnUp       int = 0
	BtnDown     int = 1
	BtnInside   int = 2
)


const (
	DirUp   int = 1
	DirIdle int = 0
	DirDown int = -1
)

type ElevatorOrder struct {
	Floor      int //1 to NumFloors
	Btn        int
	ElevatorId string
	Ack        bool
}

type ElevatorAliveMessage struct {
	Direction  int
	LastFloor  int
	ElevatorId string
}

type LightUpdate struct {
	Floor    int
	Button   int
	UpdateTo bool
}
