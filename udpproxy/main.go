package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/uzhinskiy/lib.go/pconf"
)

var (
	configfile string
	err        error
	appConfig  pconf.ConfigType
	hostname   string
)

func init() {
	flag.StringVar(&configfile, "config", "main.cfg", "Read configuration from this file")
	flag.Parse()

	appConfig = make(pconf.ConfigType)
	err := appConfig.Parse(configfile)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal("Bootstrap: error while parsing config file ", err)
	}
	hostname, _ = os.Hostname()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix(hostname + "\t")

	log.Println("Bootstrap: successful parsing config file - here is", len(appConfig), "items:", appConfig)

}

func main() {
	sconn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: appConfig["fromport"],
		IP:   net.ParseIP(appConfig["input"]),
	})
	defer sconn.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("server listening %s\n", sconn.LocalAddr().String())

	remoteAddr, err := net.ResolveUDPAddr("udp", appConfig["output"]+":"+appConfig["toport"])
	cconn, err := net.DialUDP("udp", nil, remoteAddr)
	defer cconn.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("client connected to %s\n", cconn.LocalAddr().String())

	for {
		message := make([]byte, 1500)
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
