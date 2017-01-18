package main

import(
	.fmt
	net
	../tcp
	time
)


func main(){
	String raddr="127.0.0.1:" String
	tcp.new_tcp_conn(raddr);
	Println("Connected")
}
