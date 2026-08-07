package network

import (
	"net"
	"net/netip"
)

type Client struct {
	Connection  *net.UDPConn
	EventBus    *EventBus
	PeerId      PeerId
	DoneChannel chan struct{}
}

func StartClient(address string, events *EventBus) (*Client, error) {
	addr, err := netip.ParseAddrPort(address)
	if err != nil {
		return nil, err
	}
	c, err := net.DialUDP(NETWORK_PROTOCOL, nil, net.UDPAddrFromAddrPort(addr))
	if err != nil {
		return nil, err
	}
	err = c.SetWriteBuffer(MAX_SAFE_BUFFER_SIZE)
	if err != nil {
		return nil, err
	}

	s := &Client{
		Connection:  c,
		EventBus:    events,
		PeerId:      PeerId_ConnectionRequest,
		DoneChannel: make(chan struct{}),
	}

	go s.beginListen()

	c.Write(CreatePacket(s.PeerId, nil))

	c.Write(CreatePacket(2, nil))

	return s, nil
}

func (s *Client) beginListen() {}

func (s *Client) Close() {
	s.Connection.Close()
	close(s.DoneChannel)
}
