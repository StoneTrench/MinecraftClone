package protocol

import "unsafe"

// PROTOCOL_VERSION is incremented on breaking changes.
const PROTOCOL_VERSION = 0
const NETWORK_PROTOCOL = "udp"

type PeerId uint32

const (
	// This peer has no peer id yet.
	PeerId_ConnectionRequest PeerId = 0
	// Peer id of the server.
	PeerId_Server PeerId = 1
	// Largest peer id which may not be used.
	PeerId_MAX_STATIC PeerId = PeerId_Server
	// Next largest peer id which may be used after the ones which cannot be used.
	PeerId_MIN_ALLOWED PeerId = PeerId_MAX_STATIC + 1
	// Largest peer id which may be used.
	PeerId_MAX PeerId = (1 << (SIZEOF_PeerId * 8)) - 1
)

type SeqNum uint16

const MAX_CLIENTS_ON_SERVER = int(PeerId_MAX - PeerId_MAX_STATIC)
const PROTOCOL_ID uint32 = 0x836414d1

type PayloadType uint8

const (
	PayloadType_CONTROL  PayloadType = 0
	PayloadType_ORIGINAL PayloadType = 1
	PayloadType_SPLIT    PayloadType = 2
	PayloadType_RELIABLE PayloadType = 3
)

type ControlType uint8

const (
	ControlType_SET_PEER_ID ControlType = 0
	ControlType_ACK         ControlType = 1
	ControlType_PING        ControlType = 2 // keepalive empty control packet
	ControlType_DISCONNECT  ControlType = 3
)

// Sizes of the elements

const SIZEOF_PROTOCOL_ID = int(unsafe.Sizeof(PROTOCOL_ID))
const SIZEOF_PeerId = int(unsafe.Sizeof(PeerId(0)))
const SIZEOF_SequenceNumber = int(unsafe.Sizeof(SeqNum(0)))
const SIZEOF_PayloadType = int(unsafe.Sizeof(PayloadType(0)))
const SIZEOF_ControlType = int(unsafe.Sizeof(ControlType(0)))
const SIZEOF_ChunkCount = int(unsafe.Sizeof(ChunkCount(0)))
const SIZEOF_ChunkIndex = int(unsafe.Sizeof(ChunkIndex(0)))

const HEADER_SIZE = SIZEOF_PROTOCOL_ID + SIZEOF_PeerId
const HEADER_SIZE_CONTROL = SIZEOF_PayloadType + SIZEOF_ControlType
const HEADER_SIZE_ORIGINAL = SIZEOF_PayloadType
const HEADER_SIZE_SPLIT = SIZEOF_PayloadType + SIZEOF_SequenceNumber + SIZEOF_ChunkCount + SIZEOF_ChunkIndex
const HEADER_SIZE_RELIABLE = SIZEOF_PayloadType + SIZEOF_SequenceNumber

const PACKET_SIZE = 1 << 9
const PAYLOAD_SIZE_MAX = PACKET_SIZE - HEADER_SIZE
const PAYLOAD_SIZE_ORIGINAL = PAYLOAD_SIZE_MAX - HEADER_SIZE_ORIGINAL
const PAYLOAD_SIZE_RELIABLE_ORIGINAL = PAYLOAD_SIZE_ORIGINAL - HEADER_SIZE_RELIABLE
const PAYLOAD_SIZE_SPLIT = PAYLOAD_SIZE_MAX - HEADER_SIZE_SPLIT
const PAYLOAD_SIZE_RELIABLE_SPLIT = PAYLOAD_SIZE_SPLIT - HEADER_SIZE_RELIABLE
