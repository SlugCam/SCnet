package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func listen(port string, handler func(net.Conn)) {
	go func() {
		os.Remove(port)
		ln, err := net.Listen("unix", port)
		if err != nil {
			log.Fatal(err)
		}
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Fatal("error in TCP command connection listener: ", err)
			}
			go handler(conn)
		}
	}()
}

func readLineHandler(ch chan<- string) {
	return func(c net.Conn) {
		defer c.Close()
		scanner := bufio.NewScanner(c)
		//scanner.Split(bufio.ScanLines) this is default
		for scanner.Scan() {
			line := scanner.Text()
			log.Info("message received: ", line)
			ch <- line
		}
	}
}
