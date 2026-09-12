package protocol

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

type SentReliablePacket struct {
	SentData []byte
	Retries  int
	LastSent time.Time
}
type ReceivedReliablePacket struct {
	Payload  []byte
	Received time.Time
}
type ReceivedSplitPackets struct {
	Data          [][]byte
	ChunkCount    ChunkCount
	ReceivedCount ChunkCount
	Received      time.Time
}

type ListenForPayload_Result struct {
	Payload   []byte
	PeerId    PeerId
	Address   net.Addr
	ErrClosed error
	ErrFatal  error
}

func ListenForPayload(ctx context.Context, conn net.PacketConn) <-chan ListenForPayload_Result {
	channel := make(chan ListenForPayload_Result, 10)
	buf := make([]byte, PACKET_SIZE)

	go func() {
		defer close(channel)

		var err error
		var n int
		var addr net.Addr

		for {
			select {
			case <-ctx.Done():
				channel <- ListenForPayload_Result{ErrClosed: ctx.Err()}
				return
			default:
				// Interact with networking
				err = conn.SetReadDeadline(time.Now().Add(time.Second))
				if err != nil {
					goto handle_err
				}
				n, addr, err = conn.ReadFrom(buf)
				if err != nil {
					goto handle_err
				}

				// Received
				if peer_id, payload, ok := ParsePacket_WithHeader(buf[:n]); ok {
					select {
					case channel <- ListenForPayload_Result{Payload: payload, PeerId: peer_id, Address: addr}:
					case <-ctx.Done():
					}
					return
				}

				continue

				// Handle errors
			handle_err:
				if errors.Is(err, net.ErrClosed) {
					select {
					case channel <- ListenForPayload_Result{ErrClosed: fmt.Errorf("connection was closed: %w", err)}:
					case <-ctx.Done():
					}
					return
				}
				if netErr, ok := err.(*net.OpError); ok && netErr.Timeout() {
					continue
				}

				select {
				case channel <- ListenForPayload_Result{ErrFatal: fmt.Errorf("error occured in listener: %w", err)}:
				case <-ctx.Done():
				}
				return
			}
		}
	}()

	return channel
}

type ReliableTransmitter struct {
	mu sync.Mutex

	Conn             net.PacketConn
	RemoteAddress    net.Addr
	LocalIdentifier  PeerId
	RemoteIdentifier PeerId
	LastHeardFrom    time.Time

	// Config

	RetryTimeout time.Duration
	PacketTooOld time.Duration
	MaxRetries   int

	// Send

	Buffer_SentPackets map[SeqNum]*SentReliablePacket
	Next_ReliableSeq   SeqNum
	Next_SplitSeq      SeqNum

	// Receive

	NextExpected_ReliableSeq SeqNum
	Buffer_ReliablePackets   map[SeqNum]ReceivedReliablePacket
	Buffer_SplitPackets      map[SeqNum]*ReceivedSplitPackets
	ChanIncomingData         chan []byte
}

func CreateReliableTransmitter(
	connection net.PacketConn,
	remote_address net.Addr,
	local_peer_id PeerId,
	remote_peer_id PeerId,
	retry_timeout time.Duration,
	max_retry_count int,
	packet_too_old_age time.Duration,
) *ReliableTransmitter {
	res := &ReliableTransmitter{
		mu: sync.Mutex{},

		Conn:             connection,
		RemoteAddress:    remote_address,
		LocalIdentifier:  local_peer_id,
		RemoteIdentifier: remote_peer_id,
		LastHeardFrom:    time.Now(),

		RetryTimeout: retry_timeout,
		MaxRetries:   max_retry_count,
		PacketTooOld: packet_too_old_age,

		Buffer_SentPackets: make(map[SeqNum]*SentReliablePacket),
		Next_ReliableSeq:   0,
		Next_SplitSeq:      0,

		NextExpected_ReliableSeq: 0,
		Buffer_ReliablePackets:   make(map[SeqNum]ReceivedReliablePacket),
		Buffer_SplitPackets:      make(map[SeqNum]*ReceivedSplitPackets),
		ChanIncomingData:         make(chan []byte, 10),
	}

	return res
}

// SendPayload just sends a payload without any guarantees.
func (r *ReliableTransmitter) SendPayload(payload []byte, reliable bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.send_payload(payload, reliable)
}

// HandlePacket needs to receive the listened for payloads given by ListenForPayload.
func (r *ReliableTransmitter) HandlePacket(payload []byte, addr net.Addr, remote_peer_id PeerId) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.handle_packet(payload, addr, remote_peer_id)
}

// Sends a disconnection packet, and closes the transmitter, does not close the connection.
func (r *ReliableTransmitter) Disconnect() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.send_payload(NewPayloadControl_Disconnect(), true)
	r.close()
}

// UpdateSentReliablePackets should be called on a timer.
func (r *ReliableTransmitter) UpdateSentReliablePackets() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.clean_buffer_sent_reliable()
}

// RemoveOldPackets should be called on a timer.
func (r *ReliableTransmitter) RemoveOldPackets() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clean_buffer_reliable()
}

// RemoveOldSplitPackets should be called on a timer.
func (r *ReliableTransmitter) RemoveOldSplitPackets() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clean_buffer_split_chunk()
}

// p_wrap_in_packets increments Next_SplitSeq.
func (r *ReliableTransmitter) split_payload(payload []byte, is_reliable bool) [][]byte {
	limit := PAYLOAD_SIZE_SPLIT
	if is_reliable {
		limit = PAYLOAD_SIZE_RELIABLE_SPLIT
	}
	payloads := NewPayloadSplit(payload, r.Next_SplitSeq, limit)
	if len(payloads) > 1 {
		r.Next_SplitSeq++
	}
	return payloads
}

func (r *ReliableTransmitter) send_payload(payload []byte, is_reliable bool) error {
	if r.RemoteAddress == nil {
		return fmt.Errorf("failed to send payload: remote address was nil")
	}
	payloads := r.split_payload(payload, is_reliable)
	var err error
	var wrapped [][]byte
	if is_reliable {
		wrapped, err = NewPackets_WithReliableHeaders(r.LocalIdentifier, payloads, r.Next_ReliableSeq)
	} else {
		wrapped, err = NewPackets_WithHeaders(r.LocalIdentifier, payloads)
	}
	if err != nil {
		return fmt.Errorf("failed to send payload: %w", err)
	}
	writing_errs := []error{}
	for i, data := range wrapped {
		_, err := r.Conn.WriteTo(data, r.RemoteAddress)
		if err != nil {
			writing_errs = append(writing_errs, fmt.Errorf("failed to send payload %d to %s: %w", i, r.RemoteAddress.String(), err))
		}
	}
	if is_reliable {
		for i, data := range wrapped {
			r.Buffer_SentPackets[r.Next_ReliableSeq+SeqNum(i)] = &SentReliablePacket{
				SentData: data,
				Retries:  0,
				LastSent: time.Now(),
			}
		}
		r.Next_ReliableSeq += SeqNum(len(wrapped))
	}
	return errors.Join(writing_errs...)
}

func (r *ReliableTransmitter) clean_buffer_sent_reliable() (err_out error) {
	to_delete := []SeqNum{}
	timeout_ms := r.RetryTimeout.Seconds()
	for seq, p := range r.Buffer_SentPackets {
		dur := time.Since(p.LastSent).Seconds() * math.Log(float64(p.Retries))
		if dur > timeout_ms {
			p.Retries++
			p.LastSent = time.Now()
			_, err := r.Conn.WriteTo(p.SentData, r.RemoteAddress)
			if err != nil {
				err_out = errors.Join(err_out, fmt.Errorf("failed to resend packet: %w", err))
			}
		}
		if p.Retries > r.MaxRetries {
			to_delete = append(to_delete, seq)
		}
	}
	for _, seq := range to_delete {
		delete(r.Buffer_SentPackets, seq)
	}
	return err_out
}

func (r *ReliableTransmitter) clean_buffer_reliable() {
	to_delete := []SeqNum{}
	for seq, p := range r.Buffer_ReliablePackets {
		if time.Since(p.Received) > r.PacketTooOld {
			to_delete = append(to_delete, seq)
		}
	}
	for _, seq := range to_delete {
		delete(r.Buffer_ReliablePackets, seq)
	}
}

func (r *ReliableTransmitter) clean_buffer_split_chunk() {
	to_delete := []SeqNum{}
	for seq, p := range r.Buffer_SplitPackets {
		if time.Since(p.Received) > r.PacketTooOld {
			to_delete = append(to_delete, seq)
		}
	}
	for _, seq := range to_delete {
		delete(r.Buffer_SplitPackets, seq)
	}
}

func (r *ReliableTransmitter) handle_packet(payload []byte, addr net.Addr, remote_peer_id PeerId) error {
	if remote_peer_id != r.RemoteIdentifier {
		return fmt.Errorf("invalid remote peer id, expected %d, got %d", r.RemoteIdentifier, remote_peer_id)
	}

	r.LastHeardFrom = time.Now()
	r.RemoteAddress = addr
	i := 0
	payload_type, i := GetUint[PayloadType](payload, i)

	switch payload_type {
	case PayloadType_CONTROL:
		return r.handle_packet_control(payload, i)
	case PayloadType_ORIGINAL:
		r.ChanIncomingData <- payload
		return nil
	case PayloadType_RELIABLE:
		return r.handle_packet_reliable(payload, i, addr, remote_peer_id)
	case PayloadType_SPLIT:
		return r.handle_packet_split(payload, i)
	}
	return fmt.Errorf("invalid payload type %d", payload_type)
}

func (r *ReliableTransmitter) handle_packet_control(payload []byte, i int) error {
	ctrl_type, i := GetUint[ControlType](payload, i)
	switch ctrl_type {
	case ControlType_SET_PEER_ID:
		peer_id, _ := GetUint[PeerId](payload, i)
		if r.LocalIdentifier == PeerId_ConnectionRequest {
			r.LocalIdentifier = peer_id
		}
	case ControlType_ACK:
		seq, _ := GetUint[SeqNum](payload, i)
		delete(r.Buffer_SentPackets, seq)
	case ControlType_PING: // Keepalive empty packet with no response
	case ControlType_DISCONNECT:
		r.close()
	}

	return nil
}

func (r *ReliableTransmitter) handle_packet_reliable(payload []byte, i int, addr net.Addr, remote_peer_id PeerId) error {
	seq, i := GetUint[SeqNum](payload, i)
	inner := payload[i:]

	var err_out error

	err := r.send_payload(NewPayloadControl_Ack(seq), false)
	if err != nil {
		return err
	}

	expected_seq := r.NextExpected_ReliableSeq
	if seq == expected_seq {
		r.NextExpected_ReliableSeq++
		err = r.handle_packet(inner, addr, remote_peer_id)
		err_out = errors.Join(err_out, err)
	} else {
		r.Buffer_ReliablePackets[seq] = ReceivedReliablePacket{
			Payload:  inner,
			Received: time.Now(),
		}
	}

	for {
		to_handle, to_handle_exists := r.Buffer_ReliablePackets[expected_seq]

		if to_handle_exists {
			err = r.handle_packet(to_handle.Payload, addr, remote_peer_id)
			err_out = errors.Join(err_out, err)
			r.NextExpected_ReliableSeq++
			delete(r.Buffer_ReliablePackets, expected_seq)
		} else {
			break
		}
	}

	return err_out
}

func (r *ReliableTransmitter) handle_packet_split(payload []byte, i int) error {
	seq, i := GetUint[SeqNum](payload, i)
	chunk_count, i := GetUint[ChunkCount](payload, i)
	chunk_index, i := GetUint[ChunkIndex](payload, i)
	inner := payload[i:]

	/*

		If the sequence doesnt have a buffer yet
			Create a buffer array for the split data, allocating with len chunk_count
		Set the chunkdata to the corresponding index in the buffer array

		If chunk_count == received_count then
			Assemble the split data, and emit it through the channel
		end

	*/

	slice_buffer, exists := r.Buffer_SplitPackets[seq]
	if !exists {
		slice_buffer = &ReceivedSplitPackets{
			Data:          make([][]byte, chunk_count),
			ChunkCount:    chunk_count,
			ReceivedCount: 0,
		}
		r.Buffer_SplitPackets[seq] = slice_buffer
	}

	if len(slice_buffer.Data) <= int(chunk_index) {
		return fmt.Errorf("chunk index outside the bounds of the buffer (%d >= %d)", chunk_index, len(slice_buffer.Data))
	}
	if slice_buffer.ChunkCount != chunk_count {
		return fmt.Errorf("chunk count of received split packet not equal to previous chunk size that was received (%d != %d)", chunk_count, slice_buffer.ChunkCount)
	}

	if len(slice_buffer.Data[chunk_index]) <= 0 {
		slice_buffer.ReceivedCount++
	}

	slice_buffer.Data[chunk_index] = inner
	slice_buffer.Received = time.Now()

	if slice_buffer.ReceivedCount == slice_buffer.ChunkCount {
		total_size := 0
		for _, chunk := range slice_buffer.Data {
			total_size += len(chunk)
		}
		total_data := make([]byte, total_size)
		offst := 0
		for _, chunk := range slice_buffer.Data {
			offst_end := offst + len(chunk)
			copy(total_data[offst:offst_end], chunk)
		}
		r.ChanIncomingData <- total_data
		delete(r.Buffer_SplitPackets, seq)
	}

	return nil
}

func (r *ReliableTransmitter) close() {
	r.RemoteIdentifier = PeerId_ConnectionRequest
	r.LocalIdentifier = PeerId_ConnectionRequest
	r.RemoteAddress = nil
	r.Buffer_SentPackets = nil

	r.NextExpected_ReliableSeq = 0
	r.Buffer_ReliablePackets = nil
	r.Buffer_SplitPackets = nil
	close(r.ChanIncomingData)
}
