package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type MpvEvent struct {
	Event string   `json:"event"`
	Args  []string `json:"args"`
}

// mpv <filename> --input-ipc-server <socket filepath>

func trackViewTime(socketPath string) (float64, error) {
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
			if err != io.EOF {
				return 0.0, err
			}
			continue
		}
		mpvEvent := MpvEvent{}
		err = json.Unmarshal(buf[:n], &mpvEvent)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}
		if mpvEvent.Args != nil {
			finalTimePos, err := strconv.ParseFloat(mpvEvent.Args[0], 64)
			if err != nil {
				return 0.0, err
			}
			return finalTimePos, nil
		}
	}
}
