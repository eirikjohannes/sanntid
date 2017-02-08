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

type BtnType int

const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	up     BtnType = 1
	inside BtnType = 0
	down   BtnType = -1
)

type DirType int

const (
	up   DirType = 1
	idle DirType = 0
	down DirType = -1
)

type ElevatorOrder struct {
	Floor      int //1 to NumFloors
	Btn        BtnType
	ElevatorId string
	Ack        bool
}

type ElevatorAliveMessage struct {
	Direction  DirType
	LastFloor  int
	ElevatorId string
}

type LightUpdate struct {
	Floor    int
	Button   int
	UpdateTo bool
}
