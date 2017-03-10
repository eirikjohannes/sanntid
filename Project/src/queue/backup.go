package queue

import (
	def "definitions"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// RunBackup loads backup on startup, and saves queue whenever
// there is anything on the takeBackup channel. 
//There is activity on the Takebackup channel every time a new order is placed or excecuted/removed
func RunBackup(outgoingMsg chan<- def.Message) {

	
	var backup QueueType
	backup.loadFromDisk(def.BackupFilename)
	printQueue()

	// Read last time backup was modified
	fileStat, _ := os.Stat(filename)

	// Resend all hall Orders found in backup, and add internal Orders to queue:
	for floor := 0; floor < def.NumFloors; floor++ {
		for btn := 0; btn < def.NumButtons; btn++ {
			if backup.hasOrder(floor, btn) {
				if btn == def.BtnInside {
					AddOrder(floor, btn, def.LocalElevatorId)
					// Check if time since last backup is less than OrderTimeoutDuration
				} else if !time.Now().After(fileStat.ModTime().Add(def.ElevatorOrderTimeoutDuration)) {
					AddOrder(floor, btn, def.LocalElevatorId)
				}
			}
		}
	}
	go func() {
		for {
			<-takeBackup
			log.Println(def.ColG, "Take Backup", def.ColN)
			queue.saveToDisk(filename)
		}
	}()
}

// saveToDisk saves a QueueType to disk.
func (q *QueueType) saveToDisk(filename string) {
	data, _ := json.Marshal(&q)
	ioutil.WriteFile(filename, data, 0644)
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to a QueueType
func (q *QueueType) loadFromDisk(filename string) {
	if _, err := os.Stat(filename); err == nil {
		log.Println(def.ColG, "Backup file found, processing...", def.ColN)
		data, _ := ioutil.ReadFile(filename)
		json.Unmarshal(data, q)
	}
}

func printQueue() {
	fmt.Println(def.ColB, "\n*****************************")
	fmt.Println("*       Up     Down   Inside   ")
	for f := def.NumFloors - 1; f >= 0; f-- {
		s := "* " + strconv.Itoa(f+1) + "  "
		for b := 0; b < def.NumButtons; b++ {
			if queue.hasOrder(f, b) && b != def.BtnInside {
				s += "( " + queue.Matrix[f][b].Addr[12:15] + " ) "
			} else if queue.hasOrder(f, b) {
				s += "(  x  ) "
			} else {
				s += "(     ) "
			}
		}
		fmt.Println(s)
	}
	fmt.Println("*****************************\n", def.ColN)
}
