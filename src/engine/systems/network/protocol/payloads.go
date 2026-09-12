package protocol

func NewPayloadControl_SetPeerId(peer_id PeerId) (payload []byte) {
	payload = make([]byte, HEADER_SIZE_CONTROL+SIZEOF_PeerId)
	i := 0
	i = PutUint(payload, i, PayloadType_CONTROL)
	i = PutUint(payload, i, ControlType_SET_PEER_ID)
	i = PutUint(payload, i, peer_id)
	_ = i
	return payload
}

func NewPayloadControl_Ack(seqnum SeqNum) (payload []byte) {
	payload = make([]byte, HEADER_SIZE_CONTROL+SIZEOF_SequenceNumber)
	i := 0
	i = PutUint(payload, i, PayloadType_CONTROL)
	i = PutUint(payload, i, ControlType_ACK)
	i = PutUint(payload, i, seqnum)
	_ = i
	return payload
}

func NewPayloadControl_Ping() (payload []byte) {
	payload = make([]byte, HEADER_SIZE_CONTROL)
	i := 0
	i = PutUint(payload, i, PayloadType_CONTROL)
	i = PutUint(payload, i, ControlType_PING)
	_ = i
	return payload
}

func NewPayloadControl_Disconnect() (payload []byte) {
	payload = make([]byte, HEADER_SIZE_CONTROL)
	i := 0
	i = PutUint(payload, i, PayloadType_CONTROL)
	i = PutUint(payload, i, ControlType_DISCONNECT)
	_ = i
	return payload
}

func NewPayloadOriginal(small_data []byte) (payload []byte) {
	payload = make([]byte, HEADER_SIZE_ORIGINAL+len(small_data))
	i := 0
	i = PutUint(payload, i, PayloadType_ORIGINAL)
	copy(payload[i:], small_data)
	return payload
}

type ChunkCount uint16
type ChunkIndex uint16

// NewPayloadSplit_Single creates a single split packet. Helper function for CreatePayloadSplit.
func NewPayloadSplit_Single(data []byte, seqnum SeqNum, chunk_count ChunkCount, chunk_index ChunkIndex) (payload []byte) {
	payload = make([]byte, HEADER_SIZE_SPLIT+len(data))

	i := 0
	i = PutUint(payload, i, PayloadType_SPLIT)
	i = PutUint(payload, i, seqnum)
	i = PutUint(payload, i, chunk_count)
	i = PutUint(payload, i, chunk_index)
	copy(payload[i:], data)

	return payload
}

// NewPayloadSplit splits large_data into smaller chunks and returns a list of payloads which can be reconstructed
// by the receiver.
//
// The seqnum of a split payload is different from the seqnum of a reliable packet. Here it's used to determine which chunk belongs to which data.
func NewPayloadSplit(large_data []byte, seqnum SeqNum, data_size_limit int) (payloads [][]byte) {
	if data_size_limit <= 0 || PAYLOAD_SIZE_SPLIT < data_size_limit {
		data_size_limit = PAYLOAD_SIZE_SPLIT
	}

	total_size := len(large_data)
	chunk_count := (total_size + data_size_limit - 1) / data_size_limit // ceil division

	if chunk_count == 1 {
		// return [][]byte{CreatePayloadSplit_Single(large_data, seqnum, 1, 0)}
		return [][]byte{large_data}
	}

	payloads = make([][]byte, chunk_count)
	for i := range chunk_count {
		start := i * data_size_limit
		end := min(start+data_size_limit, total_size)
		payloads[i] = NewPayloadSplit_Single(large_data[start:end], seqnum, ChunkCount(chunk_count), ChunkIndex(i))
	}

	return payloads
}

// NewPayloadReliable attaches a reliable header to the payload, this commands the receiver to reply with an ACK payload.
//
// The seqnum parameter is used to maintain ordering of received packets.
func NewPayloadReliable(in_payload []byte, seqnum SeqNum) (payload []byte) {
	payload = make([]byte, HEADER_SIZE_RELIABLE+len(in_payload))
	i := 0
	i = PutUint(payload, i, PayloadType_RELIABLE)
	i = PutUint(payload, i, seqnum)
	copy(payload[i:], in_payload)
	return payload
}
