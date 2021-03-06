package peers

import (
	def "definitions"
	"fmt"
	"net"
	"network/conn"
	"sort"
	"time"
)

var P def.PeerUpdate

const interval = def.AliveMessageInterval
const timeout = def.ElevatorTimeoutDuration

func Reciever(port int, peerUpdateCh chan<- def.PeerUpdate) {

	var buf [1024]byte

	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection

		P.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				P.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection

		P.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				P.Lost = append(P.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {

			P.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				P.Peers = append(P.Peers, k)
			}

			sort.Strings(P.Peers)
			sort.Strings(P.Lost)
			P.NumOnline = len(P.Peers)
			peerUpdateCh <- P
		}
	}
}

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}
