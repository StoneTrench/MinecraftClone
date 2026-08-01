package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

type BracketHandler struct {
	w  io.Writer
	mu sync.Mutex
	a  []slog.Attr
}

func NewBracketHandler(w io.Writer) *BracketHandler {
	return &BracketHandler{w: w}
}

func (h *BracketHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= slog.LevelDebug
}

func (h *BracketHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	t := r.Time.Format(time.DateTime)
	lvl := r.Level.String()
	msg := r.Message

	if r.NumAttrs() > 0 {
		attrs := make([]string, 0, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value.Any()))
			return true
		})
		msg = fmt.Sprintf("%s %s", msg, attrs)
	}

	line := fmt.Sprintf("[%s] [%s] %s\n", t, lvl, msg)
	_, err := h.w.Write([]byte(line))
	return err
}

func (h *BracketHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.a)+len(attrs))
	copy(newAttrs, h.a)
	copy(newAttrs[len(h.a):], attrs)
	return &BracketHandler{w: h.w, a: newAttrs}
}

func (h *BracketHandler) WithGroup(name string) slog.Handler {
	return h
}
