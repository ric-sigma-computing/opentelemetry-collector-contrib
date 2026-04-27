// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !linux

package transport // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/statsdreceiver/internal/transport"

import "net"

// setSocketBuffer sets the receive buffer size on a packet connection.
// On non-Linux platforms, SO_RCVBUFFORCE is not available, so this falls
// back to the standard SetReadBuffer (SO_RCVBUF).
func setSocketBuffer(conn net.PacketConn, size int) error {
	type readBufferSetter interface {
		SetReadBuffer(int) error
	}
	if rb, ok := conn.(readBufferSetter); ok {
		return rb.SetReadBuffer(size)
	}
	return nil
}
