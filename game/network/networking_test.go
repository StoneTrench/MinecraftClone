package network_test

import (
	"testing"
	"time"

	"github.com/StoneTrench/go-mc-clone/game/events"
	"github.com/StoneTrench/go-mc-clone/game/log"
	"github.com/StoneTrench/go-mc-clone/game/network"
)

func Test(t *testing.T) {
	log.Init(true)

	es := events.CreateEventBus()

	server, err := network.StartServer(25565, es)
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Close()

	var client *network.Client
	for range 3 {
		client, err = network.StartClient("[::]:25565", es)
		if err != nil {
			t.Error(err)
			return
		}
		defer client.Close()
	}

	time.Sleep(time.Second * 5)
}
