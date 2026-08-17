package network_test

import (
	"context"
	"testing"
	"time"

	"github.com/StoneTrench/go-mc-clone/src/engine/log"
	"github.com/StoneTrench/go-mc-clone/src/engine/network"
)

func Test(t *testing.T) {
	log.Init(true)

	log.Info("Server started")
	server, err := network.StartServer(context.Background(), 25565)
	if err != nil {
		t.Error(err)
		return
	}
	eh1 := server.PacketSource.Subscribe(func(e network.NetworkEvent) error {
		log.Infof("Server event: %v", e)
		return nil
	})
	defer server.PacketSource.Unsubscribe(eh1)
	defer server.Close()

	var client *network.Client
	for i := range 10 {
		log.Infof("Client started: %d", i)
		client, err = network.StartClient(context.Background(), "127.0.0.1:25565")
		if err != nil {
			t.Error(err)
			return
		}
		eh2 := client.PacketSource.Subscribe(func(e network.NetworkEvent) error {
			log.Infof("Client event: %v", e)
			return nil
		})
		defer client.PacketSource.Unsubscribe(eh2)
		defer client.Close()
	}

	log.Info("Started sleeping for 5s")
	time.Sleep(time.Second * 10)
}
