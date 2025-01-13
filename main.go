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

	var filePaths []string
	for _, arg := range os.Args[1:] {
		var isValid = utils.IsValidMediaFile(arg)
		fmt.Println(arg, isValid)
		if isValid {
			filePaths = append(filePaths, arg)
		}
	}

	if len(filePaths) == 0 {
		fmt.Println("No valid media files provided.")
		os.Exit(1)
	}

	args := make([]string, len(os.Args))
	copy(args, os.Args)

	var wavPaths []string
	for _, filePath := range filePaths {
		outputPath, err := utils.ConvertToWav(filePath)
		if err != nil {
			fmt.Printf("Error converting file %s: %v\n", filePath, err)
			continue
		}
		wavPaths = append(wavPaths, outputPath)

		// Replace the original file path with the WAV output path in args
		for i, arg := range args {
			if arg == filePath {
				args[i] = outputPath
			}
		}
	}

	if len(wavPaths) > 0 {
		whisperErr = utils.ExecuteWhisper(args[1:])
		if whisperErr != nil {
			fmt.Printf("Error processing files: %v\n", whisperErr)
		}
	}

	// Cleanup WAV files
	for _, outputPath := range wavPaths {
		err := os.Remove(outputPath)
		if err != nil {
			fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
		}
	}
}
