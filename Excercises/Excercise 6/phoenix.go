package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func main() {
	_, filepath1, _, _ := runtime.Caller(0)
	currentDirectory, _ := filepath.Split(filepath1)
	storageFile := currentDirectory + "/backup.dat"
	var counter = 0
	const aliveInterval = 200 * time.Millisecond
	const aliveTimeout = 1 * time.Second

	fmt.Println("''''___Starting a new process____''''''")
	backup, err := os.Open(storageFile)
	if err != nil {
		tempFile, _ := os.Create(storageFile)
		tempFile.Close()
	}

	//An infinite loop that is only exited once the file has
	//not been written to for the last aliveTimeout interval.
	for true {
		fmt.Println("Looking for new entries in file...")
		fileStatus, err := os.Stat(storageFile)
		if err != nil {
			fmt.Println("Unrecoverable error...", err.Error())
			os.Exit(0)
		}
		if time.Now().After(fileStatus.ModTime().Add(aliveTimeout)) {
			break
		}
		time.Sleep(aliveInterval)
	}

	fmt.Println("____________Starts the main process_______")

	cmd := exec.Command("gnome-terminal", "-x", "go", "run", filepath1)
	cmd.Start()

	backupReader := bufio.NewReader(backup)
	buffer, _ := backupReader.Peek(8)

	counter, _ = strconv.Atoi(string(buffer))
	backup.Close()
	backup, _ = os.Create(storageFile)

	defer func() {
		err := backup.Close()
		if err != nil {
			panic(err)
		}
	}()

	for true {
		counter++
		counterLineToBackup := strconv.Itoa(counter)
		ioutil.WriteFile(storageFile, []byte(counterLineToBackup), 0644)
		fmt.Println(counter)
		time.Sleep(aliveInterval)
	}
}
