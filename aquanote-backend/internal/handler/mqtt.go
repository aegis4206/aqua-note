package handler

import (
	"fmt"
	"log"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type SensorHook struct {
	mqtt.HookBase
}

func (h *SensorHook) ID() string {
	return "sensor-hook"
}

func (h *SensorHook) Provides(b byte) bool {
	return b == mqtt.OnPublish
}

func (h *SensorHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	// ws Broadcast
	Broadcast(pk.Payload)
	return pk, nil
}

func StartMQTTBroker() {
	broker := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	_ = broker.AddHook(new(auth.AllowHook), nil)

	tcpListener := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: ":1883",
	})
	if err := broker.AddListener(tcpListener); err != nil {
		log.Fatalf("[MQTT Broker] AddListener error: %v", err)
	}

	if err := broker.AddHook(new(SensorHook), nil); err != nil {
		log.Fatalf("[MQTT Broker] AddHook error: %v", err)
	}

	go func() {
		fmt.Println("[MQTT Broker] Listening on :1883")
		if err := broker.Serve(); err != nil {
			log.Fatalf("[MQTT Broker] Serve error: %v", err)
		}
	}()
}
