package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunCommand(name string, args ...string) []string {
	cmd := exec.Command(name, args...)

	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error running command: %v %v\n", name, args)
	}
	output := strings.Split(string(cmdOut), "\n")

	return output
}

func InteractiveCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func StreamCommand(c chan string, prefix string, name string, args ...string) {
	cmd := exec.Command(name, args...)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			c <- fmt.Sprintf("%v - %v", prefix, text)
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}

	go func() {
		defer cmdReader.Close()

		err = cmd.Wait()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
			return
		}
	}()
}
