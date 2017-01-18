// Go 1.2
// go run helloworld_go.go

package main

import (
    . "fmt"
    "runtime"
    "time"
    "sync"
)

var globalNumber int
var mlock= &sync.Mutex{}


func someGoroutine() {
    Println("Hello from a goroutine!")
}

func thread_1(){
	for i:=0; i<1000000; i++ {
		mlock.Lock()
		globalNumber++
		mlock.Unlock()
	}    
}

func thread_2(){
	for j:=0;j<10000001; j++ {
		mlock.Lock()
		globalNumber--
		mlock.Unlock()
	}
}

func main() {
    globalNumber=0
    runtime.GOMAXPROCS(runtime.NumCPU())    // I guess this is a hint to what GOMAXPROCS does...
                                            // Try doing the exercise both with and without it!
    go thread_1()                      // This spawns someGoroutine() as a goroutine
    go thread_2()
    // We have no way to wait for the completion of a goroutine (without additional syncronization of some sort)
    // We'll come back to using channels in Exercise 2. For now: Sleep.
    time.Sleep(500*time.Millisecond)
    Println("Hello from main!",globalNumber)
}
