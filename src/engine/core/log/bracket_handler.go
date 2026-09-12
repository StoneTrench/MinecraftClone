package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// https://github.com/golang/example/blob/master/slog-handler-guide/README.md#the-withgroup-method

type BracketHandler struct {
	out  io.Writer
	mu   *sync.Mutex
	goas []groupOrAttrs
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

func NewBracketHandler(w io.Writer) *BracketHandler {
	return &BracketHandler{
		out: w,
		mu:  &sync.Mutex{},
	}
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
	_, err := h.out.Write([]byte(line))
	return err
}

func (h *BracketHandler) withGroupOrAttrs(goa groupOrAttrs) *BracketHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *BracketHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *BracketHandler) WithGroup(name string) slog.Handler {
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}
