//go:build linux
// +build linux

// from https://github.com/grpc/grpc-go/internal/syscall/syscall_linux.go

package client

import (
	"fmt"
	"net"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// setTCPUserTimeout sets the TCP user timeout on a connection's socket
func setTCPUserTimeout(conn net.Conn, timeout time.Duration) error {
	tcpconn, ok := conn.(*net.TCPConn)
	if !ok {
		// not a TCP connection. exit early
		return nil
	}
	rawConn, err := tcpconn.SyscallConn()
	if err != nil {
		return fmt.Errorf("error getting raw connection: %v", err)
	}
	err = rawConn.Control(func(fd uintptr) {
		err = syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, int(timeout/time.Millisecond))
	})
	if err != nil {
		return fmt.Errorf("error setting option on socket: %v", err)
	}

	return nil
}
