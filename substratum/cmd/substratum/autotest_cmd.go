package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sgtcodfish/substratum"

	"github.com/cyrus-and/gdb"

	"github.com/mitchellh/mapstructure"
)

type GdbConnection struct {
	Conn   *gdb.Gdb
	Logger *log.Logger

	// TODO: fill these in!
	registerList    []string
	registerNumbers []string
}

type gdbRegisterNamesResponse struct {
	Payload gdbRegisterNamesResponsePayload `json:"payload",mapstructure:"payload"`
}

type gdbRegisterNamesResponsePayload struct {
	RegisterNames []string `json:"register-names",mapstructure:"register-names"`
}

func NewGDBConnection(logger *log.Logger, gdbPath string, remoteTarget string) (*GdbConnection, error) {
	conn, err := gdb.NewCmd([]string{gdbPath, "--nx", "--quiet", "--interpreter=mi2", "-ex", "set architecture riscv:rv32"}, nil)
	if err != nil {
		return nil, err
	}

	_, err = conn.CheckedSend(fmt.Sprintf("target-select remote %s", remoteTarget))
	if err != nil {
		return nil, err
	}

	resp, err := conn.CheckedSend("data-list-register-names")
	if err != nil {
		return nil, err
	}

	var registerNames gdbRegisterNamesResponse
	err = mapstructure.Decode(resp, &registerNames)
	if err != nil {
		return nil, err
	}

	registerList := substratum.GetRegisterList()
	names := resp["payload"].(map[string]interface{})["register-names"].([]interface{})
	registerNumbers := make([]string, len(registerList))

	return &GdbConnection{
		Logger: logger,
		Conn:   conn,
	}, nil
}

func (s *GdbConnection) dumpIntegerRegisters() error {
	resp, err := s.Conn.CheckedSend("data-list-register-names")
	if err != nil {
		return err
	}

	registerList := substratum.GetRegisterList()
	names := resp["payload"].(map[string]interface{})["register-names"].([]interface{})
	registerNumbers := make([]string, len(registerList))

	for i, regName := range registerList {
		for j, foundNameRaw := range names {
			foundName := foundNameRaw.(string)
			if regName == foundName {
				registerNumbers[i] = strconv.Itoa(j)
				break
			}
		}
	}

	resp, err = s.Conn.CheckedSend(fmt.Sprintf("data-list-register-values x %s", strings.Join(registerNumbers, " ")))
	if err != nil {
		return err
	}

	s.Logger.Printf("%+v", resp)

	return nil
}

func processAutotest(logger *log.Logger) error {
	gdbPath := os.Getenv("RISCV_PREFIX") + "gdb"

	gdbConn, err := NewGDBConnection(logger, gdbPath)
	if err != nil {
		return err
	}

	logger.Printf("%+v", resp)

	for i := 0; i < 1; i++ {
		resp, err = gdbConn.Conn.CheckedSend("exec-step-instruction")
		if err != nil {
			return err
		}

		logger.Printf("%+v", resp)
	}

	return nil
}
