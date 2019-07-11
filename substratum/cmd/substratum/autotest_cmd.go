package main

import (
	"log"

	"github.com/cyrus-and/gdb"
)

func processAutotest(logger *log.Logger) error {
	// TODO: properly load gdb path
	conn, err := gdb.NewCmd([]string{"riscv32-unknown-elf-gdb", "--nx", "--quiet", "--interpreter=mi2", "-ex", "set architecture riscv:rv32"}, nil)
	if err != nil {
		return err
	}

	resp, err := conn.Send("target-select remote :1234")
	if err != nil {
		return err
	}

	logger.Printf("%+v", resp)

	for i := 0; i < 100; i++ {
		resp, err = conn.Send("exec-step-instruction")
		if err != nil {
			return err
		}

		logger.Printf("%+v", resp)
	}
	return nil
}
