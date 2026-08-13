package listen

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type mpvEvent struct {
	Event string   `json:"event"`
	Args  []string `json:"args"`
}

// mpv <filename> --input-ipc-server <socket filepath> --script= <script filepath>

func TrackViewTime(socketPath string) (float64, error) {
	socketReady := false
	for range 50 { // try for up to 5 seconds
		if _, err := os.Stat(socketPath); err == nil {
			socketReady = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !socketReady {
		return 0.0, fmt.Errorf("timed out waiting for mpv socket at %s", socketPath)
	}

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
		mpvEvent := mpvEvent{}
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
