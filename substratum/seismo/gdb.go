package seismo

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

type Reg int

const (
	X0 Reg = iota
	X1
	X2
	X3
	X4
	X5
	X6
	X7
	X8
	X9
	X10
	X11
	X12
	X13
	X14
	X15
	X16
	X17
	X18
	X19
	X20
	X21
	X22
	X23
	X24
	X25
	X26
	X27
	X28
	X29
	X30
	X31
	PC
)

type GDB struct {
	conn net.Conn
}

const gdbTimeout = 5 * time.Second

func ConnectGDB(ctx context.Context, port int) (*GDB, error) {
	conn, err := waitForTCP(ctx, fmt.Sprintf("127.0.0.1:%d", port), gdbTimeout)
	if err != nil {
		return nil, err
	}
	return &GDB{conn: conn}, nil
}

func (g *GDB) Close() error {
	return g.conn.Close()
}

//
// Low-level RSP
//

func (g *GDB) send(cmd string) error {
	packet := fmt.Sprintf("$%s#%02x", cmd, checksum(cmd))
	_, err := g.conn.Write([]byte(packet))
	return err
}

func (g *GDB) recvPacket() (string, error) {
	var payload []byte
	inPacket := false
	buf := make([]byte, 1)

	for {
		_, err := g.conn.Read(buf)
		if err != nil {
			return "", err
		}

		b := buf[0]

		switch {
		case b == '+', b == '-':
			// ACK / NACK — ignore
			continue

		case b == '$':
			inPacket = true
			payload = payload[:0]

		case b == '#':
			// skip checksum bytes
			cs := make([]byte, 2)
			_, err := g.conn.Read(cs)
			if err != nil {
				return "", err
			}
			return string(payload), nil

		default:
			if inPacket {
				payload = append(payload, b)
			}
		}
	}
}

func checksum(s string) byte {
	var c byte
	for i := 0; i < len(s); i++ {
		c += s[i]
	}
	return c
}

//
// Execution control
//

func (g *GDB) Continue() error {
	return g.send("c")
}

func (g *GDB) Halt() error {
	_, err := g.conn.Write([]byte{0x03}) // Ctrl-C
	if err != nil {
		return err
	}
	// Consume stop reply
	_, err = g.recvPacket()
	return err
}

func (g *GDB) Step() error {
	if err := g.send("s"); err != nil {
		return err
	}
	_, err := g.recvPacket()
	return err
}

//
// Registers
//

func (g *GDB) ReadAllRegs() (map[Reg]uint32, error) {
	if err := g.send("g"); err != nil {
		return nil, err
	}

	payload, err := g.recvPacket()
	if err != nil {
		return nil, err
	}

	data, err := hex.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	regs := make(map[Reg]uint32)
	for i := 0; i < 33; i++ {
		regs[Reg(i)] = binary.LittleEndian.Uint32(data[i*4:])
	}
	return regs, nil
}

func (g *GDB) ReadReg(r Reg) (uint32, error) {
	if err := g.send(fmt.Sprintf("p%x", int(r))); err != nil {
		return 0, err
	}

	payload, err := g.recvPacket()
	if err != nil {
		return 0, err
	}

	data, err := hex.DecodeString(payload)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data), nil
}

//
// Memory
//

func (g *GDB) ReadMem(addr uint32, length uint32) ([]byte, error) {
	if err := g.send(fmt.Sprintf("m%X,%X", addr, length)); err != nil {
		return nil, err
	}

	payload, err := g.recvPacket()
	if err != nil {
		return nil, err
	}

	return hex.DecodeString(payload)
}
