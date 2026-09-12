package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/StoneTrench/go-mc-clone/src/engine/core/events"
	"github.com/StoneTrench/go-mc-clone/src/engine/core/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/systems/network/protocol"
)

type EventSource[T any] = events.EventSource[T]

type Server struct {
	mu           sync.RWMutex
	peer_conn    map[protocol.PeerId]*PeerClient
	next_peer_id protocol.PeerId

	listen       *net.UDPConn
	PacketSource *EventSource[NetworkEvent]
	ctx          context.Context
	ctx_cancel   context.CancelFunc
	wait_group   sync.WaitGroup
}

func StartServer(ctx context.Context, port uint16) (*Server, error) {
	addr := &net.UDPAddr{
		Port: int(port),
	}
	ln, err := net.ListenUDP(protocol.NETWORK_PROTOCOL, addr)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	s := &Server{
		mu:           sync.RWMutex{},
		peer_conn:    make(map[protocol.PeerId]*PeerClient),
		next_peer_id: protocol.PeerId_MIN_ALLOWED, // rest is handled by GetPeerId

		listen:       ln,
		PacketSource: events.NewSource[NetworkEvent](),
		ctx:          ctx,
		ctx_cancel:   cancel,
		wait_group:   sync.WaitGroup{},
	}

	s.wait_group.Add(2)
	go s.go_begin_listen()
	go s.go_begin_update()

	return s, nil
}

func (s *Server) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.listen.Close()
	s.ctx_cancel()
}

func (s *Server) go_begin_listen() {
	s.mu.RLock()
	log.Infof("Server listening on address: %s", s.listen.LocalAddr().String())
	lstnr := s.listen
	s.mu.RUnlock()
	defer s.wait_group.Done()
	var err error

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
				log.Infof("Server closed: %s", err_closed.Error())
				s.ctx_cancel()
				return
			}
			if err_fatal != nil {
				log.Errorf("Fatal error occured in server while listening: %w", err_fatal)
				return
			}

			log.Infof("Server received payload from %d: %v", peer_id, payload)

			s.mu.Lock()
			var peer *PeerClient
			if peer_id == protocol.PeerId_ConnectionRequest {
				peer_id, err = s.get_next_peer_id()
				if err != nil {
					log.Errorf("Error occured while server was connecting to new client: %w", err)
					goto mu_continue
				}
				peer = &PeerClient{
					Transmitter: protocol.CreateReliableTransmitter(
						s.listen,
						address,
						protocol.PeerId_Server,
						peer_id,
						RELIABLE_RETRY_TIMEOUT,
						MAX_RETRY_COUNT,
						TOO_OLD_TIMEOUT,
					),
				}
				s.peer_conn[peer_id] = peer
				err = peer.Transmitter.SendPayload(protocol.NewPayloadControl_SetPeerId(peer_id), false)
				if err != nil {
					log.Errorf("Error occured while server was connecting to new client: %w", err)
				}
				err = s.PacketSource.EmitImmediate(NetworkEvent{
					Id:     EventId_PeerConnected,
					Sender: peer_id,
				})
				if err != nil {
					log.Errorf("Error occured while server was connecting to new client: %w", err)
				}

			} else {
				var exists bool
				peer, exists = s.peer_conn[peer_id]
				if !exists {
					log.Errorf("Error occured in server while looking for peer: peer with id %d not exist", peer_id)
					goto mu_continue
				}
			}

			err = peer.Transmitter.HandlePacket(payload, address, peer_id)
			if err != nil {
				log.Errorf("Error occured while server was handling packet: %w", err)
			}

			HandleIncomingData(peer.Transmitter, s.PacketSource)

		mu_continue:
			s.mu.Unlock()
		}
	}
}

func (s *Server) go_begin_update() {
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

func (s *Server) update_transmitters() {
	to_delete := []protocol.PeerId{}

	// Update
	s.mu.RLock()
	for peer_id, peer := range s.peer_conn {
		if time.Since(peer.Transmitter.LastHeardFrom) > TOO_OLD_TIMEOUT {
			to_delete = append(to_delete, peer_id)
		}
		peer.Transmitter.UpdateSentReliablePackets()
		peer.Transmitter.RemoveOldPackets()
		peer.Transmitter.RemoveOldSplitPackets()
	}
	s.mu.RUnlock()

	// Delete
	if len(to_delete) > 0 {
		s.mu.Lock()
		for _, id := range to_delete {
			delete(s.peer_conn, id)
		}
		s.mu.Unlock()
	}
}

func (s *Server) get_next_peer_id() (protocol.PeerId, error) {
	conn_len := len(s.peer_conn)
	if conn_len == protocol.MAX_CLIENTS_ON_SERVER {
		return 0, fmt.Errorf("maximum possible peers reached %d == %d cannot connect more", protocol.MAX_CLIENTS_ON_SERVER, conn_len)
	}
	if conn_len > protocol.MAX_CLIENTS_ON_SERVER {
		return 0, fmt.Errorf("maximum possible peers exceeded %d < %d cannot connect more", protocol.MAX_CLIENTS_ON_SERVER, conn_len)
	}

	for {
		if s.next_peer_id < protocol.PeerId_MIN_ALLOWED {
			s.next_peer_id = protocol.PeerId_MIN_ALLOWED
		}
		_, exists := s.peer_conn[s.next_peer_id]
		if !exists {
			break
		}
		s.next_peer_id++
	}

	return s.next_peer_id, nil
}

type PeerClient struct {
	Transmitter *protocol.ReliableTransmitter
}
