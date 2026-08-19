//go:build windows

package listen

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Microsoft/go-winio"
)

func connectToPlayer(pipePath string) (net.Conn, error) {
	pipeReady := false
	for range 50 { // try for up to 5 seconds
		if _, err := os.Stat(pipePath); err == nil {
			pipeReady = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !pipeReady {
		return nil, fmt.Errorf("timed out waiting for mpv pipe at %s", pipePath)
	}

	return winio.DialPipe(pipePath, nil)
}
