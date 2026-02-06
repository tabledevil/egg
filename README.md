# CTF Egg Hunt Tool

A cool, hacky terminal-based CTF tool written in Go.

## Features
- **Hacky Aesthetics:** Retro CRT, Cyberpunk, and Minimalist themes.
- **Visual Effects:** Typewriter text, glitch transitions, and interactive decryption.
- **Obfuscated Data:** Questions and answers are baked into the binary and obfuscated to prevent cheating via .
- **Fuzzy Matching:** Lenient answer checking (ignores case, spaces, and allows typos).
- **Single Binary:** Everything is contained in one executable.

## How to Build
1. Modify `questions.json` with your own challenges.
2. Run the build script:
   ```bash
   ./build.sh
   ```
   This will pack the JSON data into the Go source and compile the binary.

## Usage
Run the tool:
```bash
./ctf-tool
```

## Customization
- **Themes:** Check `pkg/ui/theme/` to add new visual styles.
- **Logic:** Answer validation logic is in `pkg/game/logic.go`.
