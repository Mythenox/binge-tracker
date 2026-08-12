package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

// use time-pos to get playback time

type MpvEvent struct {
	Event string   `json:"event"`
	Args  []string `json:"args"`
}

func main() {
	socketPath := "/tmp/my-app.sock"

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("Unable to connect to socket: %v", err)
	}
	fmt.Printf("Connected to unix socket %s...\n", socketPath)

	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Fatalf("Connection terminated: %v", err)
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		mpvEvent := MpvEvent{}
		err = json.Unmarshal(buf[:0+n], &mpvEvent)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}
		if mpvEvent.Args != nil {
			finalTimePos, err := strconv.ParseFloat(mpvEvent.Args[0], 64)
			if err != nil {
				log.Fatalf("Error converting time pos to float: %v", err)
			}
			fmt.Printf("Final time pos of video: %f\n", finalTimePos)
			break
		}

		fmt.Printf("Received: %s\n", string(buf[:0+n]))
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			break
		}
		log.Printf("Received: %s\n", string(buf[:0+n]))
	}
}
