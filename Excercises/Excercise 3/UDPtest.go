package main

import(
	"fmt"
	"os"
	"net"
)

func checkError(err error) {
	if err != nil{
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":30000")
	checkError(err)
	

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	defer serverConn.Close()
	
	buf := make([]byte, 1024)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		fmt.Println("Recieved this shit: ", string(buf[0:n])," from: ", addr)
		
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}	


}


