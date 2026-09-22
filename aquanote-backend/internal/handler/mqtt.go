package handler

import (
	"aquanote-backend/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"aquanote-backend/internal/repository"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type SensorHook struct {
	mqtt.HookBase
	ctx         context.Context
	tempLogRepo *repository.TemperatureLogRepository
}

func (h *SensorHook) ID() string {
	return "sensor-hook"
}

func (h *SensorHook) Provides(b byte) bool {
	return b == mqtt.OnPublish
}

func (h *SensorHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	var sensorData model.SensorData
	if err := json.Unmarshal(pk.Payload, &sensorData); err != nil {
		log.Printf("[MQTT] JSON unmarshal error: %v", err)
		return pk, err
	}

	if sensorData.DeviceCode == nil || *sensorData.DeviceCode == "" {
		log.Printf("[MQTT] missing device_code, payload: %s", pk.Payload)
		return pk, fmt.Errorf("device_code is required")
	}

	savedLog, err := h.tempLogRepo.Create(h.ctx, sensorData)
	if err != nil {
		log.Printf("[DB] insert temperature log failed: %v", err)
		return pk, err
	}

	Broadcast(*savedLog)
	return pk, nil
}

func StartMQTTBroker(ctx context.Context, tempLogRepo *repository.TemperatureLogRepository) {
	broker := mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	if err := broker.AddHook(new(auth.Hook), &auth.Options{
		Ledger: &auth.Ledger{
			Auth: auth.AuthRules{
				{
					Username: auth.RString("iot-test"),
					Password: auth.RString("iot-test"),
					Allow:    true,
				},
				{Allow: false},
			},
			Users: auth.Users{
				"iot-test": auth.UserRule{
					Password: auth.RString("iot-test"),
					ACL: auth.Filters{
						"#": auth.ReadWrite,
					},
				},
			},
		},
	}); err != nil {
		log.Fatalf("[MQTT Broker] AddHook error: %v", err)
	}

	tcpListener := listeners.NewTCP(listeners.Config{
		ID:      "t1",
		Address: ":1883",
	})
	if err := broker.AddListener(tcpListener); err != nil {
		log.Fatalf("[MQTT Broker] AddListener error: %v", err)
	}

	sensorHook := &SensorHook{ctx: ctx, tempLogRepo: tempLogRepo}
	if err := broker.AddHook(sensorHook, nil); err != nil {
		log.Fatalf("[MQTT Broker] AddHook error: %v", err)
	}

	go func() {
		log.Println("[MQTT Broker] Listening on :1883")
		if err := broker.Serve(); err != nil {
			log.Fatalf("[MQTT Broker] Serve error: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("[MQTT Broker] shutting down...")
		if err := broker.Close(); err != nil {
			log.Printf("[MQTT Broker] Close error: %v", err)
		}
	}()
}
