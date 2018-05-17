package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	commandName := os.Args[1]
	options := os.Args[2:]
	if len(options) == 0 {
		fmt.Printf("executing: \"%v\"\n", commandName)
	} else {
		fmt.Printf("executing: \"%v %v\"\n", commandName, strings.Join(options[:], " "))
	}
	cmd := exec.Command(commandName, options...)

	_, _, exitCode, err := runCommand(cmd)
	fmt.Printf("exit code: %d\n", exitCode)
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(cmd *exec.Cmd) (stdout, stderr string, exitCode int, err error) {
	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	errReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	var bufout, buferr bytes.Buffer
	outReader2 := io.TeeReader(outReader, &bufout)
	errReader2 := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		return
	}

	go print(outReader2)
	go print(errReader2)

	err = cmd.Wait()

	stdout = bufout.String()
	stderr = buferr.String()

	if err != nil {
		if err2, ok := err.(*exec.ExitError); ok {
			if s, ok := err2.Sys().(syscall.WaitStatus); ok {
				err = nil
				exitCode = s.ExitStatus()
			}
		}
	}
	return
}

func print(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
}
