/*Definition file for elevator project
Group 13, Einar Henriksen, Eirik Larsen
*/

package definitions

import "time"

const NumFloors = 4
const NumButtons = 3
const ElevatorTimeoutDuration = 5 * time.Second
const AliveMessageInterval = 500 * time.Millisecond

type BtnType int 
const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	up     BtnType = 1
	inside BtnType = 0
	down   BtnType = -1
)

type ElevatorOrder struc {
	Floor int //1 to NumFloors
	Btn BtnType
	OrderId
}