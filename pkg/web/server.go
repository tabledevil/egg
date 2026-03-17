// Package web provides an embedded HTTP + WebSocket server that spawns the
// ctf-tool inside a PTY and streams it to a browser-based xterm.js client.
//
// Optimisations over a generic terminal proxy (ttyd):
//   - Frame batching: PTY output is buffered and flushed at most every 16 ms
//     (≈60 FPS on the wire), coalescing intermediate redraws.
//   - Back-pressure: if the WebSocket write buffer is congested, intermediate
//     frames are dropped rather than queued.
//   - Binary frames: raw bytes avoid per-message UTF-8 validation overhead.
//   - Zero-scrollback client: the xterm.js frontend is configured with no
//     scrollback and a WebGL renderer for minimal browser-side overhead.
package web

import (
	"context"
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
	"nhooyr.io/websocket"
)

//go:embed static/index.html
var staticFiles embed.FS

const (
	// flushInterval caps how often PTY output is forwarded to the client.
	flushInterval = 16 * time.Millisecond

	// readBufSize is the PTY read buffer size.
	readBufSize = 32 * 1024
)

// Serve starts the web terminal server on the given address (e.g. ":8080").
// selfPath is the absolute path to the running binary so we can re-exec it
// without the -web flag.
func Serve(addr string, selfPath string, extraArgs []string) error {
	mux := http.NewServeMux()

	// Serve the static HTML/JS frontend.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(data)
	})

	// WebSocket endpoint — one connection = one PTY session.
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWS(w, r, selfPath, extraArgs)
	})

	log.Printf("web terminal listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func handleWS(w http.ResponseWriter, r *http.Request, selfPath string, extraArgs []string) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// Allow any origin for local development.
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("websocket accept: %v", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Spawn the ctf-tool (ourselves without -web) inside a PTY.
	cmd := exec.CommandContext(ctx, selfPath, extraArgs...)
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
	)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Printf("pty start: %v", err)
		conn.Close(websocket.StatusInternalError, fmt.Sprintf("pty: %v", err))
		return
	}
	defer ptmx.Close()

	var wg sync.WaitGroup

	// --- PTY → WebSocket (with frame batching) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		batchedWriter(ctx, conn, ptmx)
	}()

	// --- WebSocket → PTY (input relay) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		inputRelay(ctx, conn, ptmx)
	}()

	// Wait for the child process to exit.
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.Wait()
		cancel()
	}()

	wg.Wait()
}

// batchedWriter reads from the PTY and flushes to the WebSocket at most once
// per flushInterval.  If the WebSocket write would block, intermediate data is
// dropped (the next full flush will contain the latest state).
func batchedWriter(ctx context.Context, conn *websocket.Conn, r io.Reader) {
	buf := make([]byte, readBufSize)
	var (
		mu      sync.Mutex
		pending []byte
	)

	// Flusher goroutine — sends buffered data at a fixed cadence.
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(flushInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				data := pending
				pending = nil
				mu.Unlock()

				if len(data) == 0 {
					continue
				}

				writeCtx, writeCancel := context.WithTimeout(ctx, 50*time.Millisecond)
				err := conn.Write(writeCtx, websocket.MessageBinary, data)
				writeCancel()
				if err != nil {
					// Drop frame on congestion; only fatal errors stop the loop.
					if ctx.Err() != nil {
						return
					}
				}
			}
		}
	}()

	// Reader loop — accumulates PTY output into the pending buffer.
	for {
		n, err := r.Read(buf)
		if n > 0 {
			mu.Lock()
			pending = append(pending, buf[:n]...)
			mu.Unlock()
		}
		if err != nil {
			return
		}
	}
}

// inputRelay reads messages from the WebSocket and writes them to the PTY.
// Binary messages starting with byte 1 are resize commands.
func inputRelay(ctx context.Context, conn *websocket.Conn, ptmx *os.File) {
	for {
		typ, data, err := conn.Read(ctx)
		if err != nil {
			return
		}

		if typ == websocket.MessageBinary && len(data) == 5 && data[0] == 1 {
			// Resize: [1, cols_hi, cols_lo, rows_hi, rows_lo]
			cols := binary.BigEndian.Uint16(data[1:3])
			rows := binary.BigEndian.Uint16(data[3:5])
			if cols > 0 && rows > 0 && cols < 500 && rows < 200 {
				pty.Setsize(ptmx, &pty.Winsize{
					Cols: cols,
					Rows: rows,
				})
			}
			continue
		}

		// Regular input — write to PTY.
		if len(data) > 0 {
			ptmx.Write(data)
		}
	}
}
