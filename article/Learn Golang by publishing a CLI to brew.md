## Requirements

1. Follow the general structure provided in the article below
2. Come up with a better thumbnail 
3. Populate the numbered sections with step-wise incremental code.
	1. All code changes - ie, updating existing code or adding new code - should be explained using simple comments.
4. Please follow the structure of this article as much as possible - [https://www.freecodecamp.org/news/build-a-url-shortener-in-deno/](https://www.freecodecamp.org/news/build-a-url-shortener-in-deno/). Ensure that nothing should be missing from there. Ask me if there are any questions
5. Ensure there are no grammatical or spelling mistakes. Use chatgpt for this.
6. Avoid using the words mentioned here at all costs - [https://github.com/FareedKhan-dev/Detect-AI-text-Easily/blob/main/ai_words.txt](https://github.com/FareedKhan-dev/Detect-AI-text-Easily/blob/main/ai_words.txt)
7. References to follow:
	1. https://github.com/akash-joshi/better-whisper/
	2. https://github.com/akash-joshi/better-whisper/tree/demo
	3. https://learnfastmakethings.com/p/create-a-custom-cli-tool-and-distribute-with-homebrew-using-goreleaser-and-github-actions-4c92cdeb9dfb

---
## Introduction

- Brief introduction to Better-Whisper as a CLI tool for audio transcription
- Why Golang is great for CLI tools
- What we'll build and learn

## Let's try Better-Whisper out

### Prerequisites

- Go installed - https://go.dev/
- Brew installed - https://brew.sh/

```sh
brew tap akash-joshi/homebrew-akash-joshi
# this should automatically install ffmpeg and whisper-cpp as dependencies
brew install better-whisper
```

- Need to download model locally - https://huggingface.co/ggerganov/whisper.cpp/blob/main/ggml-tiny.en-q5_1.bin - and remember path
- Download sample audio file (to be uploaded by Akash to github)
- `better-whisper -m /path/to/model -otxt /path/to/sample.mp3`
- A txt file with the transcript will be generated at the same location as the sample file

Now, let's understand how we built this ... (Expand this section to make it more wordy)

1. Exploring Existing Solutions

First we explore what already exists ...

### Using OpenAI Whisper

- Install openai-whisper using Python - `pip install -U openai-whisper`
  - Prerequisite ffmpeg - should be installed by brew install better-whisper - but can also be installed as `brew install ffmpeg`
  - https://github.com/openai/whisper
  - **Command** `whisper /path/to/sample` - this works but is very slow

## Using whisper-cpp

- https://github.com/ggerganov/whisper.cpp
- should be installed by brew install better-whisper - but can also be installed as `brew install whisper-cpp`
- **Command** `whisper-cpp -m /path/to/model -otxt /path/to/sample.mp3`
- This throws an error because it expects input in certain format
  - `ffmpeg -i input.mp3 -ar 16000 -ac 1 -c:a pcm_s16le output.wav`
- Can we automate that?

1. Go CLI Fundamentals
   - Package structure and main function
     - Start with initialising a golang project.
     - Create and run a "hello world" program on main.go

```
// all go programs start with package main
package main

import "fmt"

// main is the entry point for the program
func main() {
	fmt.Println("Hello, World!")
}
```

- Handling command line arguments
  - Create a version of the code which logs command line arguements. Run locally to check it works
- Error handling patterns in Go
  - Error is returned as second element of tuple in response.
  - No try catch needed because error handling is explicit in go

2. Working with External Commands
   - Executing system commands with os/exec

- Implement and explain how this works - `cmd := exec.Command("whisper-cpp", "-m", "/Users/efem/Documents/models/ggml-small.en-q5_1.bin", "-osrt", outputPath)`
- Pass all flags to whisper-cpp from command line

  - Managing command output and errors - Explain how this works

```
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
```

3. Audio Processing with FFmpeg

   - How to install ffmpeg-go - https://pkg.go.dev/github.com/u2takey/ffmpeg-go
   - Converting audio formats - add the `convertToWav` method to the code
   - Handling temporary files - we use `err = os.Remove(outputPath)` to delete the converted wav file after we're done with it

4. Creating a Homebrew Formula using Goreleaser
   - Steps to set up from here - https://learnfastmakethings.com/p/create-a-custom-cli-tool-and-distribute-with-homebrew-using-goreleaser-and-github-actions-4c92cdeb9dfb
   - github actions workflow - https://github.com/akash-joshi/better-whisper/blob/main/.github/workflows/release.yml
   - goreleaser yml - https://github.com/akash-joshi/better-whisper/blob/main/.goreleaser.yml

## Conclusion

- Key takeaways
- Future improvements
  - Error handling enhancements
  - Performance optimizations
- Where to learn more

## Resources

- GitHub repository
- Related tools and libraries
- Homebrew documentation
- Goreleaser documentation

