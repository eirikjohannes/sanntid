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
const ElevatorOrderTimeoutDuration = 2 * time.Second
const CostReplyTimeoutDuration = 500 * time.Millisecond

const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	BtnUp     int = 0
	BtnDown   int = 1
	BtnInside int = 2
)

const (
	DirUp   int = 1
	DirIdle int = 0
	DirDown int = -1
)

var LocalElevatorId string
var OnlineElevators int

type NumOnline int //kan slettes?

type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string
}

const ( //Category for messages
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

type ButtonPress struct {
	Floor  int
	Button int
}

type LightUpdate struct {
	Floor    int
	Button   int
	UpdateTo bool
}

type PeerUpdate struct {
	Peers     []string
	New       string
	Lost      []string
	NumOnline int
}

type MessageChan struct {
	Outgoing  chan Message
	Incoming  chan Message
	CostReply chan Message
}

type HardwareChan struct {
	MotorDir       chan int
	FloorLamp      chan int
	DoorLamp       chan bool
	BtnPressed     chan ButtonPress
	DoorTimerReset chan bool
}
type EventChan struct {
	FloorReached       chan int
	DoorTimeout        chan bool
	ElevatorPeerUpdate chan PeerUpdate
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
