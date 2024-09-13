package main

import (
	"fmt"
	"os"
	"os/exec"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func convertToWav(filePath string) (string, error) {
	outputPath := fmt.Sprintf("%s_temp.wav", filePath)

	err := ffmpeg_go.Input(filePath).
		Output(outputPath, ffmpeg_go.KwArgs{"ar": 16000, "ac": 1, "c:a": "pcm_s16le"}).
		Run()
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func executeWhisper(args []string) error {
	whisperCmd := exec.Command("whisper-cpp", args...)
	whisperCmd.Stdout = os.Stdout
	whisperCmd.Stderr = os.Stderr

	return whisperCmd.Run()
}

func main() {
	filePath := os.Args[len(os.Args)-1]

	fileExists := true
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fileExists = false
	}

	if !fileExists {
		executeWhisper(os.Args[1:])
		os.Exit(0)
	}

	outputPath, err := convertToWav(filePath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	args := append(os.Args[1:len(os.Args)-1], outputPath)
	whisperErr := executeWhisper(args)

	err = os.Remove(outputPath)
	if err != nil {
		fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
		// We don't return here, as the main operation (transcription) has already completed
	}

	if whisperErr != nil {
		fmt.Println(whisperErr)
	}
}
