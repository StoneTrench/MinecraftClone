package network

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/StoneTrench/go-mc-clone/game/log"
)

type Server struct {
	mu              sync.Mutex
	LocalUDPAddr    *net.UDPAddr
	Listener        *net.UDPConn
	EventBus        *EventBus
	PeerConnections map[PeerId]PeerClient
	DoneChannel     chan struct{}
	NextPeerId      PeerId
}

func StartServer(port uint16, events *EventBus) (*Server, error) {
	addr := &net.UDPAddr{
		Port: int(port),
	}
	ln, err := net.ListenUDP(NETWORK_PROTOCOL, addr)
	if err != nil {
		return nil, err
	}
	err = ln.SetWriteBuffer(MAX_SAFE_BUFFER_SIZE)
	if err != nil {
		return nil, err
	}

	s := &Server{
		mu:              sync.Mutex{},
		LocalUDPAddr:    addr,
		Listener:        ln,
		EventBus:        events,
		PeerConnections: make(map[PeerId]PeerClient),
		DoneChannel:     make(chan struct{}),
		NextPeerId:      PeerId_MaxStaticNext, // rest is handled by GetPeerId
	}

	go s.beginListen()

	return s, nil
}

func (s *Server) beginListen() {
	log.Infof("Server listening on address %s", s.Listener.LocalAddr().String())

	buf := make([]byte, MAX_SAFE_BUFFER_SIZE)
	for {
		err := s.Listener.SetDeadline(time.Now().Add(time.Second * 30))
		if err != nil {
			log.Errorf("Error occured while listening for new connections, cannot set deadline, %w", err)
			continue
		}
		n, addr, err := s.Listener.ReadFromUDP(buf)
		if errors.Is(err, net.ErrClosed) {
			break
		}
		if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
			continue
		}

		if err == nil && n == MAX_SAFE_BUFFER_SIZE {
			err = fmt.Errorf("packet too large, no support for packages larger than %d implemented yet", MAX_SAFE_BUFFER_SIZE)
		}
		if err != nil {
			log.Errorf("Error occured while listening for new connections, %w", err)
			continue
		}

		if peer_id, payload, ok := ParsePacket(buf); ok {
			_ = payload
			if peer_id == PeerId_ConnectionRequest {
				peer_id, err = s.GetPeerId()
				if err != nil {
					log.Errorf("Error occured while attempting to connect new peer, %w", err)
					continue
				}

				s.mu.Lock()
				s.PeerConnections[peer_id] = PeerClient{
					Address: addr,
				}
				s.mu.Unlock()

				s.EventBus.EmitDeferred(PeerConnected{
					PeerId: peer_id,
				})
				log.Infof("New peer connected with id %d at %s", peer_id, addr.String())
			} else { // Not a new connection
				s.mu.Lock()
				peer, exists := s.PeerConnections[peer_id]
				if !exists {
					log.Errorf("Received packet with invalid peer id %d, no peer with this id connected from %s", peer_id, addr.String())
					s.mu.Unlock()
					continue
				}

				is_address_valid := peer.Address.IP.Equal(addr.IP) && peer.Address.Port == addr.Port
				s.mu.Unlock()
				if !is_address_valid {
					log.Errorf("Received packet with suspicious address %s from peer %d, known address of peer should be %s", addr.String(), peer_id, peer.Address.String())
					continue
				}
			}
		}
	}
}

func (s *Server) Close() {
	s.Listener.Close()
	close(s.DoneChannel)
}

func (s *Server) GetPeerId() (PeerId, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn_len := len(s.PeerConnections)
	if conn_len == MAX_CLIENTS_ON_SERVER {
		return 0, fmt.Errorf("maximum possible peers reached %d == %d cannot connect more", MAX_CLIENTS_ON_SERVER, conn_len)
	}
	if conn_len > MAX_CLIENTS_ON_SERVER {
		return 0, fmt.Errorf("maximum possible peers exceeded %d < %d cannot connect more", MAX_CLIENTS_ON_SERVER, conn_len)
	}

	for {
		if s.NextPeerId < PeerId_MaxStaticNext {
			s.NextPeerId = PeerId_MaxStaticNext
		}
		_, exists := s.PeerConnections[s.NextPeerId]
		if !exists {
			break
		}
		s.NextPeerId++
	}

	return s.NextPeerId, nil
}

type PeerClient struct {
	Address       *net.UDPAddr
	LastHeardFrom time.Time
}
