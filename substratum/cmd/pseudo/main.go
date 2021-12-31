package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

// pty funcs based on:
// https://github.com/google/goterm/blob/master/example/example.go

/*
const (
	// TIOCGPTN is used to identify the sub pty ID
	TIOCGPTN = 0x80045430

	// TIOCSPTLCK is used to unlock the sub pty
	TIOCSPTLCK = 0x40045431
)


func getSubPtyName(masterFile *os.File) (string, error) {
	var ptyno uint64
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(masterFile.Fd()), uintptr(TIOCGPTN), uintptr(unsafe.Pointer(&ptyno)))
	if errno != 0 {
		return "", errno
	}

	name := "/dev/pts/" + strconv.FormatUint(ptyno, 10)
	return name, nil
}

func unlockSubPty(masterFile *os.File) error {
	unlock := 0 // 0 means "unlock"
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(masterFile.Fd()), uintptr(TIOCSPTLCK), uintptr(unsafe.Pointer(&unlock)))

	if errno != 0 {
		return errno
	}

	return nil
}

*/

func reader(scanner *bufio.Scanner) {
	for scanner.Scan() {
		log.Printf("qemu: %s", scanner.Text())
	}
}

// ################################# 2 ########################

func main() {
	ctx := context.Background()

	qemuArgs := []string{
		"-nographic",
		"-serial",
		"pty",
		"-s",
		"-S",
		"-M",
		"sifive_e",
		"-monitor",
		"stdio",
		"-kernel",
		"/home/ashley/workspace/bb2/06-uart-rxxd/BUILD/uart-rxxd.elf",
	}

	pipeReader, pipeWriter := io.Pipe()

	stdinPipeReader, stdinPipeWriter := io.Pipe()

	cmd := exec.CommandContext(ctx, "/usr/bin/qemu-system-riscv32", qemuArgs...)
	cmd.Stdin = stdinPipeReader
	cmd.Stdout = pipeWriter

	err := cmd.Start()
	if err != nil {
		log.Printf("failed to start qemu: %s", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(pipeReader)

	go reader(scanner)

	time.Sleep(2 * time.Second)

	_, err = stdinPipeWriter.Write([]byte("\nhelp\n"))
	if err != nil {
		log.Printf("failed to send help to qemu: %s", err)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	_, err = stdinPipeWriter.Write([]byte("\nquit\n"))
	if err != nil {
		log.Printf("failed to send quit to qemu: %s", err)
		os.Exit(1)
	}

	err = stdinPipeWriter.Close()
	if err != nil {
		log.Printf("failed to close stdin: %s", err)
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	err = cmd.Wait()
	if err != nil {
		log.Printf("failed to wait for qemu: %s", err)
		os.Exit(1)
	}

	log.Println("qemu quit successfully")
}

/*
func main2() {
	masterFile, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		log.Printf("failed to open ptmx: %s", err)
		os.Exit(1)
	}

	defer masterFile.Close()

	err = unlockSubPty(masterFile)
	if err != nil {
		log.Printf("failed to unlock sub pty: %s", err)
		os.Exit(1)
	}

	subPtyName, err := getSubPtyName(masterFile)
	if err != nil {
		log.Printf("failed to get sub pty name: %s", err)
		os.Exit(1)
	}

	log.Printf("ptyname: %s", subPtyName)

	subFile, err := os.OpenFile(subPtyName, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		log.Printf("failed to open sub pty: %s", err)
		os.Exit(1)
	}

	defer subFile.Close()

	ctx := context.Background()

	qemuArgs := []string{
		"-nographic",
		"-serial",
		subPtyName,
		"-s",
		"-S",
		"-M",
		"sifive_e",
		"-monitor",
		"stdio",
		"-kernel",
		"/home/ashley/workspace/bb2/06-uart-rxxd/BUILD/uart-rxxd.elf",
	}

	cmd := exec.CommandContext(ctx, "/usr/bin/qemu-system-riscv32", qemuArgs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = subFile, subFile, subFile
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}

	log.Printf("starting qemu")

	err = cmd.Start()
	if err != nil {
		log.Printf("failed to start qemu: %s", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(masterFile)

	go reader(scanner)

	time.Sleep(time.Second * 5)

	log.Printf("sending shutdown signal")

	n := 0
	total := 0
	n, _ = masterFile.Write([]byte("q"))
	total += n
	log.Printf("total: %d", total)
	time.Sleep(time.Millisecond * 500)
	n, _ = masterFile.Write([]byte("u"))
	total += n
	log.Printf("total: %d", total)
	time.Sleep(time.Millisecond * 500)
	n, _ = masterFile.Write([]byte("i"))
	total += n
	log.Printf("total: %d", total)
	time.Sleep(time.Millisecond * 500)
	n, _ = masterFile.Write([]byte("t"))
	total += n
	log.Printf("total: %d", total)
	time.Sleep(time.Millisecond * 500)
	n, _ = masterFile.Write([]byte("\n"))
	total += n
	log.Printf("total: %d", total)

	time.Sleep(time.Second * 5)

	err = cmd.Wait()
	if err != nil {
		log.Printf("failed to wait for qemu: %s", err)
		os.Exit(1)
	}
}
*/
