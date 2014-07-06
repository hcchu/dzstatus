package main

import (
    "bufio"
    "fmt"
    "os/exec"
)

func getTitle(tmessages chan string) {
	cmd := exec.Command("xtitle", "-s")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error")
	}
	if err := cmd.Start(); err != nil {
		fmt.Println("error")
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
        s := scanner.Text()
        tmessages <-s
    }
}
