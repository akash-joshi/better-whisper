<script>
  import { open } from "@tauri-apps/plugin-dialog";
  import { readTextFile } from "@tauri-apps/plugin-fs";

  /** @type {string | null} */
  let filePath = null;

  /** @type {string | null} */
  let srtPath = null;

  /** @type {string | null} */
  let srtContent = null;

  async function transcribe() {
    try {
      filePath = await open({
        multiple: false,
        filters: [
          {
            name: "Audio",
            extensions: ["mp3", "wav", "ogg"],
          },
        ],
      });

      if (filePath) {
        console.log("Selected file path:", filePath);
      }
    } catch (err) {
      console.error("Failed to open file:", err);
    }
  }

  async function transcribeAudio() {
    if (!filePath) return;

    try {
      const response = await fetch("http://localhost:8080/transcribe", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          filename: filePath,
          args: [
            "-m",
            "/Users/efem/Documents/models/ggml-small.en-q5_1.bin",
            "-osrt",
          ],
        }),
      });

      /** @type {{ success: boolean, srtPath: string, error?: string }} */
      const data = await response.json();

      if (data.success) {
        console.log("Transcription output path:", data.srtPath);
        srtPath = data.srtPath;
        srtContent = await readTextFile(srtPath);
      } else {
        console.error("Transcription failed:", data.error);
      }
    } catch (err) {
      console.error("API call failed:", err);
    }
  }
</script>

<main class="container">
  <button on:click={transcribe}>Select Audio File</button>
  {filePath}
  <br />
  <button on:click={transcribeAudio}>Transcribe</button>
  {#if srtContent}
    <pre>{srtContent}</pre>
  {/if}
</main>
