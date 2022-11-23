//go:build !linux
// +build !linux

package client

import (
	"net"
	"net/url"
)

func setDialControl(dialer *net.Dialer, u *url.URL) {
	// nothing to do on non-linux platforms
}
