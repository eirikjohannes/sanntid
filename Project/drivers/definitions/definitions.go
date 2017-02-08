/*Definition file for elevator project
Group 13, Einar Henriksen, Eirik Larsen
*/

package definitions

import "time"

const NumFloors = 4
const ElevatorDoorTimeoutDuration = 2 * time.Second
const ElevatorTimeoutDuration = 5 * time.Second
const AliveMessageInterval = 500 * time.Millisecond
const UDPPort = 13131

const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	up     int = 1
	inside int = 0
	down   int = -1
)

const (
	//	up   DirType = 1
	idle int = 0
	//	down DirType = -1
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
