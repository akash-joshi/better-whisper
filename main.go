package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

func whisper_print_usage() {
	fmt.Println("\nusage: whisper-cpp [options] file.wav")
	fmt.Println("\noptions:")
	fmt.Println("  -h,        --help              [default] show this help message and exit")
	fmt.Println("  -t N,      --threads N         [4      ] number of threads to use during computation")
	fmt.Println("  -p N,      --processors N      [1      ] number of processors to use during computation")
	fmt.Println("  -ot N,     --offset-t N        [0      ] time offset in milliseconds")
	fmt.Println("  -on N,     --offset-n N        [0      ] segment index offset")
	fmt.Println("  -d  N,     --duration N        [0      ] duration of audio to process in milliseconds")
	fmt.Println("  -mc N,     --max-context N     [-1     ] maximum number of text context tokens to store")
	fmt.Println("  -ml N,     --max-len N         [0      ] maximum segment length in characters")
	fmt.Println("  -sow,      --split-on-word     [false  ] split on word rather than on token")
	fmt.Println("  -bo N,     --best-of N         [5      ] number of best candidates to keep")
	fmt.Println("  -bs N,     --beam-size N       [5      ] beam size for beam search")
	fmt.Println("  -wt N,     --word-thold N      [0.01   ] word timestamp probability threshold")
	fmt.Println("  -et N,     --entropy-thold N   [2.40   ] entropy threshold for decoder fail")
	fmt.Println("  -lpt N,    --logprob-thold N   [-1.00  ] log probability threshold for decoder fail")
	fmt.Println("  -debug,    --debug-mode        [false  ] enable debug mode (eg. dump log_mel)")
	fmt.Println("  -tr,       --translate         [false  ] translate from source language to english")
	fmt.Println("  -di,       --diarize           [false  ] stereo audio diarization")
	fmt.Println("  -tdrz,     --tinydiarize       [false  ] enable tinydiarize (requires a tdrz model)")
	fmt.Println("  -nf,       --no-fallback       [false  ] do not use temperature fallback while decoding")
	fmt.Println("  -otxt,     --output-txt        [false  ] output result in a text file")
	fmt.Println("  -ovtt,     --output-vtt        [false  ] output result in a vtt file")
	fmt.Println("  -osrt,     --output-srt        [false  ] output result in a srt file")
	fmt.Println("  -olrc,     --output-lrc        [false  ] output result in a lrc file")
	fmt.Println("  -owts,     --output-words      [false  ] output script for generating karaoke video")
	fmt.Println("  -fp,       --font-path         [/System/Library/Fonts/Supplemental/Courier New Bold.ttf] path to a monospace font for karaoke video")
	fmt.Println("  -ocsv,     --output-csv        [false  ] output result in a CSV file")
	fmt.Println("  -oj,       --output-json       [false  ] output result in a JSON file")
	fmt.Println("  -ojf,      --output-json-full  [false  ] include more information in the JSON file")
	fmt.Println("  -of FNAME, --output-file FNAME [       ] output file path (without file extension)")
	fmt.Println("  -ps,       --print-special     [false  ] print special tokens")
	fmt.Println("  -pc,       --print-colors      [false  ] print colors")
	fmt.Println("  -pp,       --print-progress    [false  ] print progress")
	fmt.Println("  -nt,       --no-timestamps     [false  ] do not print timestamps")
	fmt.Println("  -l LANG,   --language LANG     [en     ] spoken language ('auto' for auto-detect)")
	fmt.Println("  -dl,       --detect-language   [false  ] exit after automatically detecting language")
	fmt.Println("             --prompt PROMPT     [       ] initial prompt")
	fmt.Println("  -m FNAME,  --model FNAME       [models/ggml-base.en.bin] model path")
	fmt.Println("  -oved D,   --ov-e-device DNAME [CPU    ] the OpenVINO device used for encode inference")
	fmt.Println("  -ls,       --log-score         [false  ] log best decoder scores of tokens")
	fmt.Println("  -ng,       --no-gpu            [false  ] disable GPU")
}

func containsHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

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

	if containsHelpFlag(os.Args) {
		whisper_print_usage()
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

	outputPath, err := convertToWav(filePath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return
	}

	args := append(os.Args[1:len(os.Args)-1], outputPath)
	whisperErr = executeWhisper(args)

	err = os.Remove(outputPath)
	if err != nil {
		fmt.Printf("Error deleting temporary file %s: %v\n", outputPath, err)
		// We don't return here, as the main operation (transcription) has already completed
	}

	if whisperErr != nil {
		fmt.Println(whisperErr)
	}
}
