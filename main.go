// main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const (
	DEBUG   = 1
	CONFIG  = "./main.conf"
	LOGFILE = "main.log"
)

type CfgVars struct {
	LogFile string
	Debug   int
}

var configfile string
var cfgvars CfgVars

func init() {
	var cfgRaw = make(map[string]string)
	flag.StringVar(&configfile, "config", CONFIG, "Read configuration from this file")
	flag.StringVar(&configfile, "c", CONFIG, "Read configuration from this file (short)")
	flag.Parse()

	rawBytes, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}

	text := string(rawBytes)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fields := strings.Split(line, "=")
		if len(fields) == 2 && strings.HasPrefix(fields[0], ";") == false {
			cfgRaw[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}
	}

	if DEBUG == 1 {
		log.Println(cfgRaw, len(cfgRaw))
	}

	if len(cfgRaw) > 0 {

		if cfgRaw["logfile"] != "" {
			cfgvars.LogFile = cfgRaw["logfile"]
		} else {
			cfgvars.LogFile = LOGFILE
		}
	}

}
func main() {
	var err error
	var i int
	var host string
	host, _ = os.Hostname()

	t := time.Now()
	time_layout := fmt.Sprintf("%02d-%02d-%02d_%02d_%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
	log_file := fmt.Sprintf(cfgvars.LogFile, time_layout)

	/* связываем вывод log-сообщений с файлом */
	logTo := os.Stderr
	if logTo, err = os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		log.Fatal(err)
	}
	defer logTo.Close()
	log.SetOutput(logTo)

	c := 0
	for {
		i++
		if c == 1202 {
			t = time.Now()
			time_layout = fmt.Sprintf("%02d-%02d-%02d_%02d_%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
			log_file = fmt.Sprintf(cfgvars.LogFile, time_layout)

			/* связываем вывод log-сообщений с файлом */
			logTo = os.Stderr
			if logTo, err = os.OpenFile(log_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
				log.Fatal(err)
			}
			defer logTo.Close()
			log.SetOutput(logTo)
			c = 0
		}
		log.Printf("%s %s[%d]:\t%d\t%s\n", host, os.Args[0], os.Getpid(), i, "some message")
		time.Sleep(999 * time.Millisecond)
		c++
	}

}
