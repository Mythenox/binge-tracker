//go:build !windows

package listen

import (
	"fmt"
	"net"
	"os"
	"time"
)

// mpv <filename> --input-ipc-server <socket filepath> --script= <script filepath>

func connectToPlayer(socketPath string) (net.Conn, error) {
	socketReady := false
	for range 50 { // try for up to 5 seconds
		if _, err := os.Stat(socketPath); err == nil {
			socketReady = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !socketReady {
		return nil, fmt.Errorf("timed out waiting for mpv socket at %s", socketPath)
	}

	return net.Dial("unix", socketPath)
}
