package main

import (
	"./network/bcast"
	"./network/localip"
	"./network/peers"
	"flag"
	"fmt"
	"os"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
}

type BtnType int

const (
	//up/down are external buttons, inside is the inside-button, floor is specified in struct newOrder
	up     BtnType = 1
	inside BtnType = 0
	down   BtnType = -1
)

type NewOrder struct {
	Floor int //1 to n
	Btn   BtnType
	Id    int
}

func main() {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)
	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(13132, id, peerTxEnable)
	go peers.Receiver(13132, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	orderTx := make(chan NewOrder)
	orderRx := make(chan NewOrder)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(13131, helloTx)
	go bcast.Transmitter(13131, orderTx)
	go bcast.Receiver(13131, helloRx)
	go bcast.Receiver(13131, orderRx)

	// The example message. We just send one of these every second.
	go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		order := NewOrder{3, 1, os.Getpid()}
		for {
			order.Floor++
			if order.Floor > 4 {
				order.Floor = 1
			}
			//fmt.Println("Floor is: %d", order.Floor)
			order.Btn--
			if order.Btn < -1 {
				order.Btn = 1
			}
			orderTx <- order
			time.Sleep(2 * time.Second)
		}
	}()
	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		case b := <-orderRx:
			fmt.Printf("NEW ORDER:\t%#v\n", b)
		}
	}
}
