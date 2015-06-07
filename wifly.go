package main

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"time"
)

const ALIVE_MESSAGE = `{"id":0,"cam":"%s","time":%d,"data":{"message":"node connected"},"type":"alive"}
`

type WiFly struct {
	serial        io.ReadWriter
	remoteAddress string
	videoPort     string
	messagePort   string
}

func NewWiFly(serial io.ReadWriter) *WiFly {
	wifly := new(WiFly)
	wifly.serial = serial
	return wifly
}

func (w *WiFly) readInput() []byte {
	time.Sleep(200 * time.Millisecond)
	return w.serial.Read()
}

func (w *WiFly) resetWiFly() {
	w.out.Write([]byte("\rreboot\r"))
	time.Sleep(3 * time.Second)
}

// Will currently lock if does not encounter correct sequence
func (w *WiFly) enterCommMode() bool {
	time.Sleep(300 * time.Millisecond)
	w.out.Write([]byte("$$$"))
	time.Sleep(300 * time.Millisecond)
	w.out.Write([]byte("\r"))

	// Clear the input
	r := w.readInput()
	if regexp.MatchString("[.\\n]*CMD", r) || regexp.MatchString("[.\\n]*<.*>", r) {
		return true
	} else {
		return false
	}
}

func (w *WiFly) openConnection(ip string, port string) bool {
	w.enterCommMode()
	w.out.Write([]byte(fmt.Sprintf("open %s %s\r", ip, port)))
	// Clear the input
	r := w.readInput
	if regexp.MatchString("[.\\n]*\\*OPEN\\*", r) {
		return true
	} else {
		w.out.Write([]byte("exit\r"))
		return false
	}
}

func (w *WiFly) checkMessages() error {
	if !openConnection(w.remoteAddress, w.messagePort) {
		return fmt.Errorf("error opening connection")
	}
	// connection is open
	// write alive message

	// read incoming messages
	// write outgoing messages

	return nil
}

func (w *WiFly) sendOneVideo() error {
	//status, err := bufio.NewReader(conn).ReadBytes('\r')

	f, err := os.Open(c.Path)
	if err != nil {
		log.Error("Gateway: Error opening video file: ", err)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		log.Error("Gateway: Error stating video file: ", err)
		return
	}

	log.WithFields(log.Fields{
		"path": c.Path,
	}).Info("Gateway sending video to central server")

	// Camera name
	fmt.Fprintf(conn, "%d\x00", c.Source)

	// Timestamp/ID
	id := make([]byte, 4)
	binary.LittleEndian.PutUint32(id, uint32(c.Timestamp.Unix()))
	conn.Write(id)

	// Data size
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(fi.Size()))
	conn.Write(size)

	// Send data
	n, err := io.Copy(conn, f)
	if n != fi.Size() || err != nil {
		log.WithFields(log.Fields{
			"written": n,
			"total":   fi.Size(),
			"error":   err,
		}).Error("Gateway: Error writing video data")
	}

	return nil
}

func (w *WiFly) handleWiFly(mch <-chan string, vch <-chan string) {
	for {
		//checkMessages()
		sendOneVideo()
	}
}
