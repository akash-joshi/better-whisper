# Better-Whisper 🎙️ - Config-free Whisper Transcription

Better-Whisper is a command-line interface tool that uses the Whisper speech recognition model, providing easy audio file conversion and transcription capabilities. It handles all media formats so you can transcribe any file without configuration.

## Quick Start

You can quickly run the Better-Whisper directly from GitHub using the `go run` command:

```sh
go run github.com/akash-joshi/better-whisper@v0.1.1 [whisper-cpp arguments] <input-file>
```

## Example

```sh
better-whisper -m ~/Documents/ggml-model-whisper-small.en.bin -t 4 -p 2 -ml 21 -sow -osrt
```

## Features

- Converts various audio formats to WAV for Whisper processing
- Executes Whisper transcription on audio files
- Handles both direct Whisper commands and audio file inputs

## Pre-requisites

You need to have [`ffmpeg`](https://formulae.brew.sh/formula/ffmpeg) and [`whisper-cpp`](https://formulae.brew.sh/formula/whisper-cpp) installed on your system.

## Usage

Instant usage:

```sh
go run main.go [whisper-cpp arguments] <input-file>
```

Build the project:

```sh
go build -o better-whisper main.go
```

Run the CLI tool:

```sh
./better-whisper [whisper-cpp arguments] <input-file>
```

If the input file exists, it will be converted to WAV format (if necessary) before being passed to Whisper. If the input file doesn't exist, the command will be passed directly to Whisper.

## How it works

1. The tool checks if the provided file exists.
2. If the file exists, it's converted to a 16kHz, mono, 16-bit PCM WAV file using FFmpeg.
3. The Whisper model is then executed on the converted file.
4. After transcription, the temporary WAV file is deleted.

## Error Handling

Errors during file conversion or Whisper execution are printed to the console. The tool attempts to clean up temporary files even if errors occur during transcription.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.