//go:build linux
// +build linux

package client

import (
	"net"
	"net/url"
	"syscall"

	"golang.org/x/sys/unix"
)

const tcpUserTimeoutMilliseconds = 20000

func setDialControl(dialer *net.Dialer, u *url.URL) {
	if u.Scheme == TCP || u.Scheme == SSL {
		dialer.Control = func(network, address string, c syscall.RawConn) error {
			var syscallErr error
			controlErr := c.Control(func(fd uintptr) {
				syscallErr = syscall.SetsockoptInt(
					int(fd), syscall.IPPROTO_TCP, unix.TCP_USER_TIMEOUT, tcpUserTimeoutMilliseconds)
			})
			if syscallErr != nil {
				return syscallErr
			}
			if controlErr != nil {
				return controlErr
			}
			return nil
		}
	}
}
