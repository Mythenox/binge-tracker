package listen

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/Microsoft/go-winio"
)

type mpvEvent struct {
	Event string   `json:"event"`
	Args  []string `json:"args"`
}

type vlcEvent struct {
	Event   string  `json:"event"`
	TimePos float64 `json:"time_seconds"`
}

// mpv <filename> --input-ipc-server <socket filepath> --script= <script filepath>

func TrackViewTimeMPVLinux(socketPath string) (float64, error) {
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

func TrackViewTimeMPVWin(pipePath string) {
	fmt.Printf("Connecting to mpv pipe at %s...\n", pipePath)

	// Set a timeout duration for the initial connection
	timeout := 10 * time.Second

	// Connect to the Windows Named Pipe
	winio.DialPipe
	conn, err := winio.DialPipe(pipePath, &timeout)
	if err != nil {
		log.Fatalf("Error connecting to pipe: %v. Is mpv running?", err)
	}
	defer conn.Close()

	fmt.Println("Connected! Listening for the video stop JSON payload...")
}

// VLC <filename> --extraintf luaintf --lua-intf <script filepath>
// --play-and-exit maybe
// --start-time

func TrackViewTimeVLC(playerCmd *exec.Cmd) (float64, error) {
	// listen to tcp port then start player

	listener, err := net.Listen("tcp", ":9099")
	if err != nil {
		return 0.0, fmt.Errorf("Failed to bind port: %v", err)
	}
	defer listener.Close()

	err = playerCmd.Start()
	if err != nil {
		return 0.0, fmt.Errorf("error launching video player: %v", err)
	}

	go func() {
		err := playerCmd.Wait()
		if err != nil {
			log.Printf("Process exited or was killed by user: %v", err)
			return
		}
	}()

	conn, err := listener.Accept()
	if err != nil {
		return 0.0, err
	}

	defer conn.Close()

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return 0.0, err
	}

	vlcEvent := vlcEvent{}
	err = json.Unmarshal(buf[:n], &vlcEvent)
	if err != nil {
		return 0.0, fmt.Errorf("Error unmarshalling event: %v", err)
	}

	finalTimePos := vlcEvent.TimePos

	return finalTimePos, nil
}
