// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build linux

package transport // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/transport"

import (
	"fmt"
	"net"
	"syscall"
)

// setSocketBuffer sets the receive buffer size on a packet connection.
// It tries SO_RCVBUFFORCE first (bypasses net.core.rmem_max, requires
// CAP_NET_ADMIN), falling back to SO_RCVBUF if unprivileged.
func setSocketBuffer(conn net.PacketConn, size int) error {
	sc, ok := conn.(syscall.Conn)
	if !ok {
		return fmt.Errorf("connection does not support syscall.Conn")
	}
	raw, err := sc.SyscallConn()
	if err != nil {
		return fmt.Errorf("getting raw conn: %w", err)
	}
	var setsockoptErr error
	err = raw.Control(func(fd uintptr) {
		setsockoptErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUFFORCE, size)
		if setsockoptErr != nil {
			// Fall back to regular SO_RCVBUF (capped by rmem_max).
			setsockoptErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_RCVBUF, size)
		}
	})
	if err != nil {
		return err
	}
	return setsockoptErr
}
