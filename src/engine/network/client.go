package network

import (
	"context"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/StoneTrench/go-mc-clone/src/engine/events"
	"github.com/StoneTrench/go-mc-clone/src/engine/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/network/protocol"
)

type Client struct {
	mu          sync.RWMutex
	Transmitter *protocol.ReliableTransmitter

	listen       *net.UDPConn
	PacketSource *EventSource[NetworkEvent]
	ctx          context.Context
	ctx_cancel   context.CancelFunc
	wait_group   sync.WaitGroup
}

func StartClient(ctx context.Context, address string) (*Client, error) {
	addr, err := netip.ParseAddrPort(address)
	if err != nil {
		return nil, err
	}
	remote_addr := net.UDPAddrFromAddrPort(addr)
	ln, err := net.ListenUDP(protocol.NETWORK_PROTOCOL, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := &Client{
		mu: sync.RWMutex{},
		Transmitter: protocol.CreateReliableTransmitter(
			ln,
			remote_addr,
			protocol.PeerId_ConnectionRequest,
			protocol.PeerId_Server,
			RELIABLE_RETRY_TIMEOUT,
			MAX_RETRY_COUNT,
			TOO_OLD_TIMEOUT,
		),

		listen:       ln,
		PacketSource: events.NewSource[NetworkEvent](),
		ctx:          ctx,
		ctx_cancel:   cancel,
		wait_group:   sync.WaitGroup{},
	}

	s.wait_group.Add(2)
	go s.go_begin_listen()
	go s.go_begin_update()

	err = s.Transmitter.SendPayload(protocol.NewPayloadOriginal(nil), true)
	if err != nil {
		log.Errorf("Error occured in client: %w", err)
	}

	return s, nil
}

func (s *Client) go_begin_listen() {
	s.mu.RLock()
	log.Infof("Client listening on address: %s", s.listen.LocalAddr().String())
	lstnr := s.listen
	s.mu.RUnlock()
	defer s.wait_group.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case listen_res := <-protocol.ListenForPayload(s.ctx, lstnr):
			payload := listen_res.Payload
			address := listen_res.Address
			peer_id := listen_res.PeerId
			err_closed := listen_res.ErrClosed
			err_fatal := listen_res.ErrFatal

			if err_closed != nil {
				log.Infof("Client closed: %s", err_closed.Error())
				s.ctx_cancel()
				return
			}
			if err_fatal != nil {
				log.Errorf("Fatal error occured in client while listening: %w", err_fatal)
				return
			}

			log.Infof("Client received payload from %d: %v", peer_id, payload)

			if peer_id == protocol.PeerId_Server {
				s.mu.RLock()
				trans := s.Transmitter
				event_bus := s.PacketSource
				s.mu.RUnlock()
				err := trans.HandlePacket(payload, address, peer_id)
				if err != nil {
					log.Errorf("Error occured while client was handling packet: %w", err)
				}

				HandleIncomingData(trans, event_bus)
			}
		}
	}
}

func (s *Client) go_begin_update() {
	defer s.wait_group.Done()
	timer := time.NewTimer(RELIABLE_RETRY_TIMEOUT)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			s.update_transmitters()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Client) update_transmitters() {
	s.mu.RLock()
	s.Transmitter.UpdateSentReliablePackets()
	s.Transmitter.RemoveOldPackets()
	s.Transmitter.RemoveOldSplitPackets()
	s.mu.RUnlock()
}

func (s *Client) Close() {
	s.Transmitter.Disconnect()
}
