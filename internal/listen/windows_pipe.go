//go:build windows

package listen

import (
	"fmt"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

func connectToPlayer(pipePath string) (net.Conn, error) {
	fmt.Printf("Connecting to mpv pipe at %s...\n", pipePath)

	// Set a timeout duration for the initial connection
	timeout := 10 * time.Second

	return winio.DialPipe(pipePath, &timeout)
}
