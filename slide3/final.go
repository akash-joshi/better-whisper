package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

// define structs for json formatting
type TranscriptionResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	SrtPath string `json:"srtPath,omitempty"`
}

type TranscriptionRequest struct {
	Filename string   `json:"filename"`
	Args     []string `json:"args"`
}

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
	cmd := exec.Command("whisper-cpp", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func handleTranscription(w http.ResponseWriter, r *http.Request) {
	var req TranscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(TranscriptionResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}

	outputPath, err := convertToWav(req.Filename)
	if err != nil {
		json.NewEncoder(w).Encode(TranscriptionResponse{
			Success: false,
			Error:   fmt.Sprintf("Error converting file: %v", err),
		})
		return
	}

	srtPath := outputPath + ".srt"

	args := append(req.Args, outputPath)
	whisperErr := executeWhisper(args)

	err = os.Remove(outputPath)
	if err != nil {
		fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
	}

	if whisperErr != nil {
		json.NewEncoder(w).Encode(TranscriptionResponse{
			Success: false,
			Error:   whisperErr.Error(),
		})
		return
	}

	// json encoding
	json.NewEncoder(w).Encode(TranscriptionResponse{
		Success: true,
		SrtPath: srtPath,
	})
}

func main() {
	// gorilla mux for routing
	r := mux.NewRouter()
	r.HandleFunc("/transcribe", handleTranscription).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1420", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
