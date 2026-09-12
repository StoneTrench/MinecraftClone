package network

import (
	"time"

	"github.com/StoneTrench/go-mc-clone/src/engine/core/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/systems/network/protocol"
)

type EventId uint8

const (
	EventId_PeerConnected EventId = iota
	EventId_PacketReceived
)

type NetworkEvent struct {
	Id     EventId
	Sender protocol.PeerId
	Data   []byte
}

const RELIABLE_RETRY_TIMEOUT = 2 * time.Second
const TOO_OLD_TIMEOUT = 30 * time.Second
const MAX_RETRY_COUNT = 10

func HandleIncomingData(t *protocol.ReliableTransmitter, e *EventSource[NetworkEvent]) {
	for {
		select {
		case d := <-t.ChanIncomingData:
			if len(d) > 0 {
				err := e.EmitImmediate(NetworkEvent{
					Id:     EventId_PacketReceived,
					Sender: t.RemoteIdentifier,
					Data:   d,
				})
				if err != nil {
					log.Errorf("Error occured while emitting network event: %w", err)
				}
			}
		default:
			return
		}
	}
}
