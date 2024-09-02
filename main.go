package main

import (
	"fmt"
	"os"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as an argument.")
		return
	}

	filePath := os.Args[1]
	outputPath := fmt.Sprintf("%s_temp.wav", filePath)

	err := ffmpeg_go.Input(filePath).
		Output(outputPath, ffmpeg_go.KwArgs{"ar": 16000, "ac": 1, "c:a": "pcm_s16le"}).
		Run()
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	fmt.Println("File converted successfully:", outputPath)
}
