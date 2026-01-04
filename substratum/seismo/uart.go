package seismo

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"
)

type UART struct {
	Conn net.Conn
	R    *bufio.Reader
}

func ConnectUART(ctx context.Context, port int, timeout time.Duration) (*UART, error) {
	addr := net.JoinHostPort("127.0.0.1", fmt.Sprint(port))

	conn, err := waitForTCP(ctx, addr, timeout)
	if err != nil {
		return nil, err
	}

	return &UART{
		Conn: conn,
		R:    bufio.NewReader(conn),
	}, nil
}

func (u *UART) ReadByte(timeout time.Duration) (byte, error) {
	_ = u.Conn.SetReadDeadline(time.Now().Add(timeout))
	return u.R.ReadByte()
}

func (u *UART) Write(b []byte) error {
	_, err := u.Conn.Write(b)
	return err
}

func (u *UART) Close() error {
	return u.Conn.Close()
}
