package main

import (
	"flag"
	"log"

	"github.com/tarm/serial"
)

const (
	MESSAGE_PORT = "/tmp/scnet_m.str"
	VIDEO_PORT   = "/tmp/scnet_v.str"
	BUFF_SIZE    = 256
)

var (
	address     string
	messagePort string
	videoPort   string
)

func main() {

	addressFlag := flag.String("a", "128.114.59.16", "server ip")
	mPort := flag.String("mp", "7892", "port for the message server")
	vPort := flag.String("vp", "7893", "port for the video server")
	baudRate := flag.Int("baud", 115200, "the baud rate for the serial connection")
	serialDev := flag.String("serial", "/dev/ttyAMA0", "path of the serial device to use")
	flag.Parse()

	address = *addressFlag
	messagePort = *mPort
	videoPort = *vPort

	// Setup serial
	c := &serial.Config{Name: *serialDev, Baud: *baudRate}
	serial, err := serial.OpenPort(c)
	if err != nil {
		log.Panic(err)
	}

	mch := make(chan string, BUFF_SIZE)
	vch := make(chan string, BUFF_SIZE)

	listen(MESSAGE_PORT, readLineHandler(mch))
	listen(VIDEO_PORT, readLineHandler(vch))

	w := &WiFly{
		serial:        serial,
		remoteAddress: address,
		messagePort:   messagePort,
		videoPort:     videoPort,
	}
	w.handleWiFly(mch, vch)
}
