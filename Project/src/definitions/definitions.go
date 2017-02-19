/*Definition file for elevator project
Group 13, Einar Henriksen, Eirik Larsen
*/

package definitions

import "time"

const NumFloors = 4
const NumButtons = 3
const ElevatorDoorTimeoutDuration = 2 * time.Second
const ElevatorTimeoutDuration = 500 * time.Millisecond
const AliveMessageInterval = 50 * time.Millisecond
const UDPPort = 13131

var LocalElevatorId string

type NumOnline int

const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	up     int = 1
	inside int = 0
	down   int = -1
	idle int = 0
)

type Message struct{
	Category	int
	Floor	int
	Button	int
	Cost	int
	Addr	string
}

const {//Category for messages
	Alive int= iota+1
	NewRequest
	CompleteRequest
	Cost
}

type ButtonPress struct{
	Floor int
	Button int
} 

type LightUpdate struct {
	Floor    int
	Button   int
	UpdateTo bool
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type MessageChan struct {
	Outgoing 	chan Message
	Incoming 	chan Message
	CostReply 	chan Message
	NumOnline 	chan int
}

type HardwareChan struct {
	MotorDir       chan int
	FloorLamp      chan int
	DoorLamp       chan bool
	BtnPressed     chan BtnPress
	DoorTimerReset chan bool
}
type EventChan struct {
	FloorReached chan int
	DoorTimeout  chan bool
	DeadElevator chan int
}



// Colors for printing to console
const Col0 = "\x1b[30;1m" // Dark grey
const ColR = "\x1b[31;1m" // Red
const ColG = "\x1b[32;1m" // Green
const ColY = "\x1b[33;1m" // Yellow
const ColB = "\x1b[34;1m" // Blue
const ColM = "\x1b[35;1m" // Magenta
const ColC = "\x1b[36;1m" // Cyan
const ColW = "\x1b[37;1m" // White
const ColN = "\x1b[0m"    // Grey (neutral)
