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