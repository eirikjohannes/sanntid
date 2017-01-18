package main

import(
	"fmt"
	"net"
	"bufio"
	"strings"
)

func checkError(err error) {
	if err != nil{
		fmt.Println("Error: ", err)
//		os.Exit(0)
	}
}

func main() {
	/*serverAddr, err := net.ResolveUDPAddr("udp", ":30000")
	checkError(err)
	

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)
	defer serverConn.Close()
	*/
	fmt.Println("Launching særver")
	
	link, err := net.Listen("tcp", ":20013")
	checkError(err)
	
	conn, err :=link.Accept()

	
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message recieved: ", string(message))

		newmessage := strings.ToUpper(message)
		
		conn.Write([]byte(newmessage + "\n"))		
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}	


}


