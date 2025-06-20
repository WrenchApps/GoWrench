package connections

import (
	"log"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/connection_settings"

	"github.com/nats-io/nats.go"
)

var connections map[string]*nats.Conn
var jetStreams map[string]nats.JetStreamContext

func GetNatsConnectionById(natsConnectionId string) *nats.Conn {
	if len(natsConnectionId) == 0 || connections == nil {
		return nil
	}

	return connections[natsConnectionId]
}

func GetJetStreamByConnectionId(natsConnectionId string) nats.JetStreamContext {
	if len(natsConnectionId) == 0 || connections == nil {
		return nil
	}

	return jetStreams[natsConnectionId]
}

func loadConnectionNats(connNatsSetting []*connection_settings.NatsConnectionSettings) error {
	if len(connNatsSetting) > 0 {
		if connections == nil {
			connections = make(map[string]*nats.Conn)
		}

		for _, conn := range connNatsSetting {
			nc, err := nats.Connect(conn.ServerAddress)

			if err != nil {
				log.Printf("Error nats connection: %v", err)
				return err
			}

			connections[conn.Id] = nc
		}
	}

	return nil
}

func loadJetStreams(settings []*action_settings.ActionSettings) error {
	var err error
	if len(settings) > 0 {
		for _, setting := range settings {
			if setting.Type == action_settings.ActionTypeNatsPublish {
				if setting.Nats.IsStream {
					if jetStreams == nil {
						jetStreams = make(map[string]nats.JetStreamContext)
					}

					if jetStreams[setting.Nats.ConnectionId] != nil {
						continue
					}

					conn := GetNatsConnectionById(setting.Nats.ConnectionId)
					js, err := conn.JetStream()

					if err != nil {
						log.Printf("Error nats jetstream: %v", err)
						return err
					}

					jetStreams[setting.Nats.ConnectionId] = js
				}
			}
		}
	}

	return err
}
