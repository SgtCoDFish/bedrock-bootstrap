package seismo

import (
	"context"
	"fmt"
	"net"
	"time"
)

// waitForTCP retries until timeout
func waitForTCP(ctx context.Context, addr string, timeout time.Duration) (net.Conn, error) {
	deadline := time.Now().Add(timeout)

	dialer := net.Dialer{}

	for time.Now().Before(deadline) {
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			return conn, nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return nil, fmt.Errorf("timed out waiting for connection to %s", addr)
}
