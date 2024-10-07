To create a Homebrew command (also known as a Homebrew formula) for your codebase, you'll need to follow these steps:

1. Create a Formula file
2. Build and package your Go application
3. Write the Formula
4. Test the Formula locally
5. Submit the Formula to Homebrew (optional)

Here's a step-by-step guide:

1. Create a Formula file:
   Create a new Ruby file with the name of your command, e.g., `better-whisper.rb`, in a local Homebrew tap directory or in the Homebrew core formulae directory.

2. Build and package your Go application:
   Before writing the formula, you need to build your Go application and create a tarball or zip file of the binary. You can do this by:

   ```sh
   go build -o better-whisper
   tar czf better-whisper-v0.1.0.tar.gz better-whisper
   ```

   Then, upload this tarball to a publicly accessible URL (e.g., a GitHub release).

3. Write the Formula:
   Here's a basic template for your `better-whisper.rb` formula:

   ```ruby
   class BetterWhisper < Formula
     desc "CLI wrapper for Whisper speech recognition"
     homepage "https://github.com/akash-joshi/better-whisper"
     url "https://github.com/akash-joshi/better-whisper/archive/v0.1.0.tar.gz"
     sha256 "replace_with_actual_sha256_of_your_tarball"
     license "MIT"

     depends_on "go" => :build
     depends_on "ffmpeg"
     depends_on "whisper-cpp"

     def install
       system "go", "build", "-o", bin/"better-whisper"
     end

     test do
       assert_match "usage: whisper-cpp", shell_output("#{bin}/better-whisper -h", 2)
     end
   end
   ```

   Replace the URL with the actual URL of your tarball, and update the SHA256 hash accordingly. You can get the SHA256 hash by running `shasum -a 256 better-whisper-v0.1.0.tar.gz`.

4. Test the Formula locally:
   You can test your formula locally by running:

   ```sh
   brew install --build-from-source ./better-whisper.rb
   ```

   This will install your formula from the local file.

5. Submit the Formula to Homebrew (optional):
   If you want to make your formula available to others via Homebrew, you'll need to submit it to the appropriate Homebrew tap or to Homebrew core. This usually involves creating a pull request on GitHub.

Remember to update your README.md with instructions on how to install your tool via Homebrew once the formula is available.

Note: The exact process might vary depending on your specific setup and requirements. You may need to adjust the formula based on your project's structure and dependencies.