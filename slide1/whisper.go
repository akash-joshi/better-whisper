package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("whisper-cpp", os.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	whisperErr := cmd.Run()
	if whisperErr != nil {
		fmt.Println(whisperErr)
	}
}
