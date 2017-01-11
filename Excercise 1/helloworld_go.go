// Go 1.2
// go run helloworld_go.go

package main

import (
    . "fmt"
    "runtime"
    "time"
)

var globalNumber int


func someGoroutine() {
    Println("Hello from a goroutine!")
}

func thread_1(){
	for i:=0; i<1000000; i++ {
		globalNumber++
	}    
}

func thread_2(){
	for j:=0;j<1000000; j++ {
		globalNumber--
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
    time.Sleep(100*time.Millisecond)
    Println("Hello from main!",globalNumber)
}
