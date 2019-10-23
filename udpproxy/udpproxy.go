package main

import (
	"fmt"
	"net"
	"log"
)

func main() {
	sconn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 514,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer sconn.Close()
	fmt.Printf("server listening %s\n", sconn.LocalAddr().String())


	RemoteAddr, err := net.ResolveUDPAddr("udp", "10.4.3.229:1514")
        cconn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer cconn.Close()
	fmt.Printf("client connected to %s\n", cconn.LocalAddr().String())



	for {
		message := make([]byte, 1024)
		rlen, _, err := sconn.ReadFromUDP(message[:])
		if err != nil {
			log.Println(err)
		}
		_, err = cconn.Write(message[:rlen])

		if err != nil {
            	    log.Println(err)
		}
		
		
	}
}