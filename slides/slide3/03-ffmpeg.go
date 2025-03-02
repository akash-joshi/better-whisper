package main

import (
	"fmt"
	"os"
	"os/exec"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func convertToWav(filePath string) (string, error) {
	outputPath := fmt.Sprintf("%s.wav", filePath)

	err := ffmpeg_go.Input(filePath).
		Output(outputPath, ffmpeg_go.KwArgs{"ar": 16000, "ac": 1, "c:a": "pcm_s16le"}).
		Run()
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func main() {
	filePath := os.Args[len(os.Args)-1]

	outputPath, err := convertToWav(filePath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	fmt.Println("Converted file to WAV:", outputPath)

	cmd := exec.Command("whisper-cpp", os.Args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	whisperErr := cmd.Run()
	if whisperErr != nil {
		fmt.Println(whisperErr)
	}
}
