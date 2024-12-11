package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/akash-joshi/better-whisper/utils"
)

var version = "devel"

func main() {
	// Check if ffmpeg exists
	_, ffmpegErr := exec.LookPath("ffmpeg")
	_, whisperErr := exec.LookPath("whisper-cpp")
	missingDependencyError := ""

	if ffmpegErr != nil {
		missingDependencyError += "\033[1mffmpeg\033[0m is missing. Please install it using Homebrew - https://formulae.brew.sh/formula/ffmpeg.\n"
	}

	// Check if whisper-cpp exists
	if whisperErr != nil {
		missingDependencyError += "\033[1mwhisper-cpp\033[0m is missing. Please install it using Homebrew - https://formulae.brew.sh/formula/whisper-cpp.\n"
	}

	if missingDependencyError != "" {
		fmt.Println(missingDependencyError)
		os.Exit(1)
	}

	if utils.ContainsVersionFlag(os.Args) {
		fmt.Println("better-whisper ", version)
		os.Exit(0)
	}

	if utils.ContainsHelpFlag(os.Args) {
		utils.WhisperPrintUsage()
		os.Exit(0)
	}

	filePath := os.Args[len(os.Args)-1]

	fileExists := true
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fileExists = false
	} else {
		if len(os.Args) == 1 || filePath == "" {
			fileExists = false
		}
	}

	if !fileExists {
		fmt.Println("No file provided or file does not exist.")
		os.Exit(1)
	}

	outputPath, err := utils.ConvertToWav(filePath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	args := append(os.Args[1:len(os.Args)-1], outputPath)
	whisperErr = utils.ExecuteWhisper(args)

	err = os.Remove(outputPath)
	if err != nil {
		fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
		// We don't return here, as the main operation (transcription) has already completed
	}

	if whisperErr != nil {
		fmt.Println(whisperErr)
	}
}
