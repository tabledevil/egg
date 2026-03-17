# Performance Optimization Plan for ttyd/Web Terminal

## Problem Statement
When hosted inside ttyd (WebSocket-based web terminal), the 30 FPS animated themes cause stutter and input lag. The root causes are:
1. **Full-screen redraws every 33ms** — every frame regenerates the entire W×H grid of ANSI-styled cells, even unchanged ones
2. **Per-cell lipgloss.Render()** — `canvas.Render()` calls `cell.Style.Render(string(cell.Rune))` for every single cell (e.g. 80×24 = 1,920 calls/frame), each producing redundant ANSI reset/set sequences
3. **No ANSI coalescing** — adjacent cells with identical styles emit separate escape sequences instead of being batched into one styled run
4. **800-particle physics computed unconditionally** — Antigravity theme runs trig and physics for all 800 particles every tick, even off-screen ones
5. **`lipgloss.NewStyle()` allocations on every Clear/New** — canvas initialization creates a fresh Style struct per cell
6. **`fmt.Sprintf` for every animated-color particle** — hex color strings allocated per-particle per-frame in AnimatedRGB/Fire/Sparkle themes
7. **ttyd overhead** — generic WebSocket terminal adds encoding/transport overhead on top of already-large ANSI payloads

## Plan (Two Tracks)

---

### Track A: Optimize Terminal Rendering (Go-side)

These changes reduce ANSI output volume and CPU cost without changing any visuals.

#### A1. Canvas: ANSI Run-Length Coalescing in `Render()`
**File:** `pkg/ui/canvas/canvas.go`

Replace the current per-cell `cell.Style.Render(string(cell.Rune))` loop with a run-coalescing renderer:
- Walk each row left-to-right, accumulate consecutive cells that share the same `lipgloss.Style` into a single string run
- Render the entire run with one `style.Render(accumulatedChars)` call
- This reduces ANSI escape sequences by ~60-80% (most backgrounds are uniform)

```go
func (c *Canvas) Render() string {
    var b strings.Builder
    for y := 0; y < c.Height; y++ {
        var run strings.Builder
        currentStyle := c.Grid[y][0].Style
        run.WriteRune(c.Grid[y][0].Rune)
        for x := 1; x < c.Width; x++ {
            cell := c.Grid[y][x]
            if cell.Style == currentStyle {  // need style equality check
                run.WriteRune(cell.Rune)
            } else {
                b.WriteString(currentStyle.Render(run.String()))
                run.Reset()
                run.WriteRune(cell.Rune)
                currentStyle = cell.Style
            }
        }
        b.WriteString(currentStyle.Render(run.String()))
        if y < c.Height-1 {
            b.WriteRune('\n')
        }
    }
    return b.String()
}
```

**Note:** lipgloss.Style is a value type. We need a reliable equality check. We'll add a `StyleKey` (compact hash/string of the style's foreground+background+attrs) to each Cell for fast comparison, since lipgloss.Style doesn't expose an equality operator.

#### A2. Canvas: Cache a default "blank" style
**File:** `pkg/ui/canvas/canvas.go`

Create a package-level `var defaultStyle = lipgloss.NewStyle()` and reuse it in `New()` and `Clear()` instead of calling `lipgloss.NewStyle()` per cell (saves ~1,920 allocations per frame).

#### A3. Antigravity: Skip off-screen particles
**File:** `pkg/ui/theme/antigravity.go`

In the particle update+draw loop, skip physics calculations for particles whose `(ix, iy)` are outside `[0, width) × [0, height)`. Currently physics runs for all 800 even if they wrapped out of bounds temporarily. Add early-continue before the avoidance/collision checks.

#### A4. Antigravity: Pre-compute color styles for animated themes
**File:** `pkg/ui/theme/antigravity.go`

For `ColorThemeAnimatedRGB`, `ColorThemeFire`, and `ColorThemeSparkle`:
- Build a small lookup table of ~64 pre-computed styles (e.g. 8×8 grid for RGB, 64 heat levels for Fire) at theme init and on each tick
- Index into the table instead of calling `fmt.Sprintf("#%02x%02x%02x", r, g, b)` + `lipgloss.NewStyle().Foreground()` for each of 800 particles
- This eliminates ~800 `fmt.Sprintf` + style allocation calls per frame

#### A5. Diff-based frame output (Bubble Tea layer)
**File:** `pkg/ui/model.go`

Cache the previous frame's rendered string in the Model. In `View()`, if the new output equals the cached output, return the cached string (Bubble Tea already does some diffing, but this avoids the string comparison overhead at the lipgloss layer). More importantly: Bubble Tea's renderer already diffs lines — our job is just to make the strings it receives shorter/more cacheable by doing A1-A4.

#### A6. Reduce tick rate for static/low-animation themes
**File:** `pkg/ui/theme/base.go`, individual themes

Add an optional `TickInterval() time.Duration` method to the Theme interface (with a default of 33ms). Themes like C64, MS-DOS, Amiga that only animate a blinking cursor can use 100ms ticks (10 FPS) instead of 30 FPS — reducing work by 3×. The model reads this and adjusts its tick command.

---

### Track B: Custom Web Terminal Server

Build a lightweight, purpose-built web terminal that replaces ttyd for this specific use case. It serves the Go app's PTY output over WebSocket with optimizations tuned for our animated TUI.

#### B1. Embedded HTTP+WebSocket server
**File:** `pkg/web/server.go` (new)

- Add a `-web` flag and optional `-port` (default 8080) to `main.go`
- When `-web` is set, instead of running BubbleTea directly on stdout:
  1. Start an HTTP server on the specified port
  2. Serve a single-page HTML terminal client (embedded via `embed`)
  3. On WebSocket connection, spawn the ctf-tool in a PTY (`os.StartProcess` or `creack/pty`)
  4. Relay PTY output → WebSocket, WebSocket input → PTY
- Dependencies: `creack/pty` for PTY management, `gorilla/websocket` or `nhooyr.io/websocket` for WebSocket

#### B2. Terminal frontend (xterm.js)
**File:** `pkg/web/static/index.html` (new, embedded)

- Single HTML file with embedded xterm.js (from CDN or vendored)
- Connect to `ws://host:port/ws`
- Configure xterm.js with performance-tuned settings:
  - `rendererType: 'canvas'` (or WebGL renderer for GPU acceleration)
  - `fastScrollModifier: 'alt'`
  - `scrollback: 0` (no scrollback needed — alt screen mode)
  - Custom font settings optimized for the retro aesthetic
- Handle resize: send terminal dimensions back over WebSocket on window resize
- Mobile-friendly viewport meta tag

#### B3. Server-side output optimization
**File:** `pkg/web/server.go`

Apply server-side optimizations before sending data over WebSocket:
- **Frame batching**: Buffer PTY output and flush at most every 16ms (60 FPS cap on the wire), coalescing intermediate frames that the user would never see
- **Bandwidth tracking**: Monitor WebSocket send buffer; if it backs up, skip frames (drop non-latest frames) rather than queuing them
- **Binary WebSocket frames**: Send raw bytes instead of text frames to avoid UTF-8 validation overhead

#### B4. Input optimization
**File:** `pkg/web/server.go`

- Debounce/batch input from WebSocket → PTY at a reasonable rate
- Send resize events as a structured message (not inline escape sequences)

#### B5. Build integration
**File:** `build.sh`

- Embed the static HTML/JS/CSS into the Go binary using `//go:embed`
- No separate build step needed — `go build` produces a single binary that can serve itself as a web terminal

---

### Implementation Order

1. **A2** (trivial, instant win — cached default style)
2. **A1** (biggest single optimization — ANSI coalescing)
3. **A4** (color LUT for Antigravity — removes per-particle allocations)
4. **A3** (off-screen particle skip)
5. **A6** (adaptive tick rate for calm themes)
6. **A5** (frame caching in model)
7. **B1** (web server skeleton + PTY)
8. **B2** (xterm.js frontend)
9. **B3** (frame batching/dropping)
10. **B4** (input handling)
11. **B5** (embed + build)

### What Does NOT Change
- No themes removed or visually altered
- No transitions removed
- No animations simplified or frame rates reduced for dynamic themes
- No game logic changes
- Existing CLI mode (`./ctf-tool`) works exactly as before
- Web mode is opt-in via `-web` flag
