package transport

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

// StdioTransport implements Transport using standard input/output.
type StdioTransport struct {
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	scanner *bufio.Scanner
	msgCh   chan []byte
	errCh   chan error
	closeCh chan struct{}
	mu      sync.Mutex
	closed  bool
}

// NewStdioTransport creates a new StdioTransport using os.Stdin and os.Stdout.
func NewStdioTransport() *StdioTransport {
	return NewStdioTransportWithIO(os.Stdin, os.Stdout, os.Stderr)
}

// NewStdioTransportWithIO creates a StdioTransport with custom IO streams.
// This is useful for testing.
func NewStdioTransportWithIO(stdin io.Reader, stdout, stderr io.Writer) *StdioTransport {
	scanner := bufio.NewScanner(stdin)
	// Increase buffer size for large JSON messages (1MB)
	const maxScanTokenSize = 1024 * 1024
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	return &StdioTransport{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		scanner: scanner,
		msgCh:   make(chan []byte, 10),
		errCh:   make(chan error, 1),
		closeCh: make(chan struct{}),
	}
}

// Start begins reading from stdin in a background goroutine.
func (t *StdioTransport) Start(ctx context.Context) error {
	go t.readLoop(ctx)
	return nil
}

// readLoop continuously reads lines from stdin and sends them to the message channel.
func (t *StdioTransport) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			t.errCh <- ctx.Err()
			return
		case <-t.closeCh:
			return
		default:
			if t.scanner.Scan() {
				line := t.scanner.Bytes()
				if len(line) == 0 {
					continue // Skip empty lines
				}
				// Make a copy since scanner reuses the buffer
				msg := make([]byte, len(line))
				copy(msg, line)
				select {
				case t.msgCh <- msg:
				case <-t.closeCh:
					return
				case <-ctx.Done():
					return
				}
			} else {
				if err := t.scanner.Err(); err != nil {
					t.errCh <- fmt.Errorf("scanner error: %w", err)
				} else {
					t.errCh <- io.EOF
				}
				return
			}
		}
	}
}

// Send writes a message to stdout followed by a newline.
func (t *StdioTransport) Send(message []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return fmt.Errorf("transport is closed")
	}

	// Write message followed by newline
	if _, err := t.stdout.Write(message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if _, err := t.stdout.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// Receive waits for and returns the next message from stdin.
func (t *StdioTransport) Receive() ([]byte, error) {
	msg, err := t.ReceiveWithMetadata()
	if err != nil {
		return nil, err
	}
	return msg.Data, nil
}

// ReceiveWithMetadata waits for and returns the next message with metadata.
// For stdio transport, metadata is always empty.
func (t *StdioTransport) ReceiveWithMetadata() (*Message, error) {
	select {
	case data := <-t.msgCh:
		return &Message{Data: data, Metadata: make(map[string]string)}, nil
	case err := <-t.errCh:
		return nil, err
	case <-t.closeCh:
		return nil, fmt.Errorf("transport is closed")
	}
}

// Close shuts down the transport.
func (t *StdioTransport) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.closed {
		return nil
	}

	t.closed = true
	close(t.closeCh)
	return nil
}

// Log writes a message to stderr for debugging purposes.
// This is separate from the main transport to avoid interfering with JSON-RPC.
func (t *StdioTransport) Log(format string, args ...interface{}) {
	fmt.Fprintf(t.stderr, format+"\n", args...)
}
