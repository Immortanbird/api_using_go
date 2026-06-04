package utils

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type ThrottledReader struct {
	reader  io.Reader       // The original source (e.g., the open file)
	limiter *rate.Limiter   // The token bucket
	ctx     context.Context // Context to handle client disconnection
}

// NewThrottledReader creates a wrapper that limits reading speed.
// limit: bytes per second
func NewThrottledReader(r io.Reader, limit int, ctx context.Context) *ThrottledReader {
	// rate.Limit is defined as "events per second". Here, 1 event = 1 byte.
	// The second parameter is "burst". We usually set burst equal to the limit
	// or a small multiple of it to allow immediate start.
	limiter := rate.NewLimiter(rate.Limit(limit), limit)

	return &ThrottledReader{
		reader:  r,
		limiter: limiter,
		ctx:     ctx,
	}
}

// Read implements the standard io.Reader interface.
// Before reading from the file, it asks the limiter to wait.
func (t *ThrottledReader) Read(p []byte) (n int, err error) {
	burst := t.limiter.Burst()
	if burst < len(p) {
		p = p[:burst]
	} else {
		burst = len(p)
	}

	// WaitN blocks (sleeps) until tokens are available or the context is canceled.
	err = t.limiter.WaitN(t.ctx, burst)
	if err != nil {
		// If the context is canceled (user disconnects), return the error.
		return 0, err
	}

	return t.reader.Read(p)
}
