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

func main() {
	// initialisations start with :=. Type inference is used here.
	filePath := "./temp/sample.mp3"

	// go methods return tuples of value, error
	outputPath, err := convertToWav(filePath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	cmd := exec.Command("whisper-cpp", "-m", "/Users/efem/Documents/models/ggml-small.en-q5_1.bin", "-osrt", outputPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	whisperErr := cmd.Run()
	if whisperErr != nil {
		fmt.Println(whisperErr)
	}

	err = os.Remove(outputPath)

	// nil is the only empty value in go
	if err != nil {
		fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
		// We don't return here, as the main operation (transcription) has already completed
	}
}
