package protocol

import (
	"fmt"
	"unsafe"
)

/*

All numbers are big-endian.

Q: What is this supposed to do?
A: It is the layer that exists between the Server/Client architecture and UDP. It handles reliable packets, split packets etc.





*/

// ParsePacket_WithHeader returns ok as true if it parsed the packet,
// or false if it didn't parse any packet.
func ParsePacket_WithHeader(buf []byte) (peer_id PeerId, payload []byte, ok bool) {
	prot_id, i := GetUint[uint32](buf, 0)
	if len(buf) < HEADER_SIZE || prot_id != PROTOCOL_ID {
		return 0, nil, false
	}
	peer_id, i = GetUint[PeerId](buf, i)
	return peer_id, buf[i:], true
}
func NewPacket_WithHeader(peer_id PeerId, payload []byte) ([]byte, error) {
	payload_size := len(payload)
	if payload_size > PAYLOAD_SIZE_MAX {
		return nil, fmt.Errorf("failed to add packet header: payload cannot be larger than %d bytes, actual size was %d bytes from peer %d", PAYLOAD_SIZE_MAX, payload_size, peer_id)
	}
	packet := NewPacket_WithHeader_Truncate(peer_id, payload)
	return packet, nil
}
func NewPacket_WithHeader_Truncate(peer_id PeerId, payload []byte) []byte {
	payload_size := min(len(payload), PAYLOAD_SIZE_MAX)
	packet := make([]byte, HEADER_SIZE+payload_size)
	i := 0
	i = PutUint(packet, i, PROTOCOL_ID)
	i = PutUint(packet, i, peer_id)
	copy(packet[i:], payload[:payload_size])
	return packet
}

func PutUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](b []byte, offset int, value T) int {
	switch unsafe.Sizeof(value) {
	case 1:
		b[offset] = byte(value)
		return offset + 1
	case 2:
		_ = b[offset+1]
		e := uint16(value)
		b[offset] = byte(e >> 8)
		b[offset+1] = byte(e)
		return offset + 2
	case 4:
		_ = b[offset+3]
		e := uint32(value)
		b[offset] = byte(e >> 24)
		b[offset+1] = byte(e >> 16)
		b[offset+2] = byte(e >> 8)
		b[offset+3] = byte(e)
		return offset + 4
	case 8:
		_ = b[offset+7]
		e := uint64(value)
		b[offset] = byte(e >> 56)
		b[offset+1] = byte(e >> 48)
		b[offset+2] = byte(e >> 40)
		b[offset+3] = byte(e >> 32)
		b[offset+4] = byte(e >> 24)
		b[offset+5] = byte(e >> 16)
		b[offset+6] = byte(e >> 8)
		b[offset+7] = byte(e)
		return offset + 8
	}
	return 0
}
func GetUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](b []byte, offset int) (T, int) {
	var value T
	switch unsafe.Sizeof(value) {
	case 1:
		return T(b[offset]), offset + 1
	case 2:
		_ = b[offset+1]
		return T(
			uint16(b[offset])<<8 |
				uint16(b[offset+1]),
		), offset + 2
	case 4:
		_ = b[offset+3]
		return T(
			uint32(b[offset])<<24 |
				uint32(b[offset+1])<<16 |
				uint32(b[offset+2])<<8 |
				uint32(b[offset+3]),
		), offset + 4
	case 8:
		_ = b[offset+7]
		return T(
			uint64(b[offset])<<56 |
				uint64(b[offset+1])<<48 |
				uint64(b[offset+2])<<40 |
				uint64(b[offset+3])<<32 |
				uint64(b[offset+4])<<24 |
				uint64(b[offset+5])<<16 |
				uint64(b[offset+6])<<8 |
				uint64(b[offset+7]),
		), offset + 8
	}
	panic("invalid type in GetUint")
}

func NewPackets_WithHeaders(local_peer_id PeerId, payloads [][]byte) ([][]byte, error) {
	packet_arr := make([][]byte, len(payloads))
	for i, p := range payloads {
		p, err := NewPacket_WithHeader(local_peer_id, p)
		if err != nil {
			return nil, fmt.Errorf("failed to construct packet: %w", err)
		}
		packet_arr[i] = p
	}
	return packet_arr, nil
}
func NewPackets_WithReliableHeaders(local_peer_id PeerId, payloads [][]byte, reliable_seq SeqNum) ([][]byte, error) {
	packet_arr := make([][]byte, len(payloads))
	for i, p := range payloads {
		p, err := NewPacket_WithHeader(local_peer_id, NewPayloadReliable(p, reliable_seq+SeqNum(i)))
		if err != nil {
			return nil, fmt.Errorf("failed to construct reliable packet: %w", err)
		}
		packet_arr[i] = p
	}
	return packet_arr, nil
}
