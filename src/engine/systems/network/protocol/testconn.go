package protocol

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type VirtualPacket struct {
	from net.Addr
	to   net.Addr
	data []byte
}

type VirtualRouter struct {
	mu    sync.RWMutex
	conns map[string]chan VirtualPacket
	done  chan struct{}
}

func NewRouter() *VirtualRouter { return &VirtualRouter{conns: map[string]chan VirtualPacket{}, done: make(chan struct{})} }

func (r *VirtualRouter) register(addr net.Addr, ch chan VirtualPacket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-r.done:
		return net.ErrClosed
	default:
	}
	if _, ok := r.conns[addr.String()]; ok {
		return fmt.Errorf("addr %s registered", addr)
	}
	r.conns[addr.String()] = ch
	return nil
}

func (r *VirtualRouter) unregister(addr net.Addr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, addr.String())
}

func (r *VirtualRouter) route(pkt VirtualPacket) error {
	r.mu.RLock()
	ch, ok := r.conns[pkt.to.String()]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no route to %s", pkt.to)
	}
	select {
	case <-r.done:
		return net.ErrClosed
	case ch <- pkt:
		return nil
	}
}

func (r *VirtualRouter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-r.done:
		return net.ErrClosed
	default:
		close(r.done)
		for _, ch := range r.conns {
			close(ch)
		}
		r.conns = nil
	}
	return nil
}

type VirtualConnection struct {
	router        *VirtualRouter
	addr          net.Addr
	recv          chan VirtualPacket
	done          chan struct{}
	closed        bool
	mu            sync.RWMutex
	ReadDeadline  time.Time
	WriteDeadline time.Time
}

func NewConn(r *VirtualRouter, addr net.Addr) (*VirtualConnection, error) {
	ch := make(chan VirtualPacket, 10)
	if err := r.register(addr, ch); err != nil {
		return nil, err
	}
	return &VirtualConnection{router: r, addr: addr, recv: ch, done: make(chan struct{})}, nil
}

func (c *VirtualConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return net.ErrClosed
	}
	c.closed = true
	c.router.unregister(c.addr)
	close(c.recv)
	close(c.done)
	return nil
}
func (c *VirtualConnection) LocalAddr() net.Addr { return c.addr }
func (c *VirtualConnection) SetDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ReadDeadline, c.WriteDeadline = t, t
	return nil
}
func (c *VirtualConnection) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ReadDeadline = t
	return nil
}
func (c *VirtualConnection) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.WriteDeadline = t
	return nil
}

func (c *VirtualConnection) ReadFrom(buf []byte) (int, net.Addr, error) {
	c.mu.RLock()
	deadline := c.ReadDeadline
	c.mu.RUnlock()

	var ctx context.Context
	var cancel context.CancelFunc
	if deadline.IsZero() {
		ctx = context.Background()
		cancel = func() {}
	} else {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()

	select {
	case <-c.done:
		return 0, nil, net.ErrClosed
	case <-ctx.Done():
		return 0, nil, &net.OpError{
			Op:   "read",
			Net:  "udp",
			Addr: c.addr,
			Err:  context.DeadlineExceeded,
		}
	case pkt, ok := <-c.recv:
		if !ok {
			return 0, nil, net.ErrClosed
		}
		n := min(len(pkt.data), len(buf))
		copy(buf[:n], pkt.data[:n])
		return n, pkt.from, nil
	}
}

func (c *VirtualConnection) WriteTo(p []byte, addr net.Addr) (int, error) {
	c.mu.RLock()
	deadline := c.WriteDeadline
	closed := c.closed
	c.mu.RUnlock()
	if closed {
		return 0, net.ErrClosed
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if deadline.IsZero() {
		ctx = context.Background()
		cancel = func() {}
	} else {
		ctx, cancel = context.WithDeadline(context.Background(), deadline)
	}
	defer cancel()

	pkt := VirtualPacket{from: c.addr, to: addr, data: p}
	errCh := make(chan error, 1)
	go func() { errCh <- c.router.route(pkt) }()

	select {
	case <-c.done:
		return 0, net.ErrClosed
	case <-ctx.Done():
		return 0, &net.OpError{
			Op:   "write",
			Net:  "udp",
			Addr: c.addr,
			Err:  context.DeadlineExceeded,
		}
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
		return len(p), nil
	}
}
