package httputil

import (
	"net"
	"strconv"
	"sync"
)

var portUseMutex sync.Mutex

func portInUseForVersion(version string, port int) bool {
	var host string

	if version == "tcp4" {
		host = "127.0.0.1:" + strconv.Itoa(port)
	} else {
		// Concatenate a colon and the port
		host = "[::1]:" + strconv.Itoa(port)
	}

	// Attempt to listen on the port
	listener, err := net.Listen(version, host)

	// If there was an error, the port is in use
	if err != nil {
		return true
	}

	// Close the listener
	listener.Close()

	return false
}

func PortInUse(port int) bool {
	// Lock the mutex
	portUseMutex.Lock()
	defer portUseMutex.Unlock()

	return portInUseForVersion("tcp6", port) ||
		portInUseForVersion("tcp4", port)
}
