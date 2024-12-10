package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

func ConvertToWav(filePath string) (string, error) {
	outputPath := fmt.Sprintf("%s_temp.wav", filePath)

	err := ffmpeg_go.Input(filePath).
		Output(outputPath, ffmpeg_go.KwArgs{"ar": 16000, "ac": 1, "c:a": "pcm_s16le"}).
		Run()
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func ExecuteWhisper(args []string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		brewPrefix, err := exec.Command("brew", "--prefix", "whisper-cpp").Output()
		if err != nil {
			return fmt.Errorf("error getting whisper-cpp prefix: %v", err)
		}
		metalPath := strings.TrimSpace(string(brewPrefix)) + "/share/whisper-cpp"

		os.Setenv("GGML_METAL_PATH_RESOURCES", metalPath)
		cmd = exec.Command("whisper-cpp", args...)
	} else {
		cmd = exec.Command("whisper-cpp", args...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func WhisperPrintUsage() {
	fmt.Println("\nusage: whisper-cpp [options] file.wav")
	fmt.Println("\noptions:")
	fmt.Println("  -h,        --help              [default] show this help message and exit")
	// ... rest of the usage print statements ...
}

func ContainsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func ContainsVersionFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			return true
		}
	}
	return false
}

func GetModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return info.Main.Version
}
