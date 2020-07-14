package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/tarm/serial"
	"github.com/utahta/go-cronowriter"
)

var (
	port       = flag.String("c", "", "serial port")
	baud       = flag.Int("b", 115200, "baud rate")
	logLoc     = flag.String("l", "", "log location")
	logPrefix  = flag.String("p", "seraildump.log", "log prefix")
	isSendTest = flag.Bool("t", false, "sent test")
)

func bgWrite(s *serial.Port) {
	for {
		time.Sleep(time.Second)
		now := time.Now()
		_, err := s.Write([]byte(now.String() + "\r\n"))
		// log.Println(now.String())
		if err != nil {
			log.Fatal("writefatal", err)
		}
	}
}

func main() {
	flag.Parse()

	logPath := path.Join(*logLoc, *logPrefix)
	fmt.Println("logging to", logPath)
	w := cronowriter.MustNew(logPath+".%Y%m%d",
		cronowriter.WithSymlink(logPath),
		cronowriter.WithInit(),
	)
	defer w.Close()

	c := &serial.Config{Name: *port, Baud: *baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal("OpenPort", err)
	}

	if *isSendTest {
		go bgWrite(s)
	}

	sc := bufio.NewScanner(s)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			ts := time.Now().Format("2006-01-02 15:04:05")
			w.Write([]byte(ts + " " + l + "\n"))
		}
	}
}
