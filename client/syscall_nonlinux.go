//go:build !linux
// +build !linux

package client

import (
	"net"
	"time"
)

// setTCPUserTimeout is a no-op function under non-linux environments
func setTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	return nil
}
