package protocol

import (
	"context"
	"net"
	"net/netip"
	"sync"
	"testing"
	"time"
)

func Test(t *testing.T) {
	waiter := sync.WaitGroup{}

	receive := func(conn net.PacketConn, trans *ReliableTransmitter) {
		defer waiter.Done()
		for {
			res := <-ListenForPayload(context.Background(), conn)
			if res.ErrClosed != nil {
				t.Errorf("Received closed error: %v", res.ErrClosed)
				break
			}
			if res.ErrFatal != nil {
				t.Errorf("Received fatal error: %v", res.ErrFatal)
				break
			}
			t.Logf("Received packet from %s: %v", res.Address.String(), res.Payload)
			err := trans.HandlePacket(res.Payload, res.Address, res.PeerId)
			if err != nil {
				t.Error(err)
			}
		}
	}
	update := func(trans *ReliableTransmitter) {
		defer waiter.Done()
		timer := time.NewTimer(time.Second)
		defer timer.Stop()
		for range timer.C {
			trans.RemoveOldPackets()
			trans.RemoveOldSplitPackets()
			trans.UpdateSentReliablePackets()
		}
	}

	t.Log("Starting virtual network")

	router := NewRouter()

	server_addr := net.UDPAddrFromAddrPort(netip.MustParseAddrPort("10.0.0.2:9001"))
	client_addr := net.UDPAddrFromAddrPort(netip.MustParseAddrPort("10.0.0.2:9002"))

	server_conn, _ := NewConn(router, server_addr)
	client_conn, _ := NewConn(router, client_addr)

	client_trans := CreateReliableTransmitter(
		client_conn,
		server_addr,
		PeerId_ConnectionRequest,
		PeerId_Server,
		1*time.Second,
		10,
		30*time.Second,
	)
	server_trans := CreateReliableTransmitter(
		server_conn,
		client_addr,
		PeerId_Server,
		PeerId_ConnectionRequest,
		1*time.Second,
		10,
		30*time.Second,
	)

	waiter.Add(2)
	go receive(server_conn, server_trans)
	go update(server_trans)

	waiter.Add(2)
	go receive(client_conn, client_trans)
	go update(client_trans)

	t.Log("Sleeping for 1s")
	time.Sleep(time.Second * 1)

	t.Log("Sending payload")
	for range 2 {
		err := client_trans.SendPayload(NewPayloadOriginal(make([]byte, 1024)), true)
		if err != nil {
			t.Error(err)
		}
	}

	waiter.Wait()
}
