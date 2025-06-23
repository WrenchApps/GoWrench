package connections

import (
	"context"
	"crypto/tls"
	"errors"
	"time"
	"wrench/app/manifest/connection_settings"

	"log"

	"github.com/segmentio/kafka-go"
)

var kafkaConnections map[string]*kafka.Conn

func loadConnectionsKafka(kafkaSettings []*connection_settings.KafkaConnectionSettings) error {

	if len(kafkaSettings) > 0 && kafkaConnections == nil {
		kafkaConnections = make(map[string]*kafka.Conn)
	}

	for _, kafkaSetting := range kafkaSettings {
		network := getKafkaNetworking(kafkaSetting.ConnectionType)

		var dialer *kafka.Dialer
		if kafkaSetting.ConnectionType == connection_settings.KafkaConnectionSsl {
			dialer = &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
				TLS:       &tls.Config{InsecureSkipVerify: true},
			}
		} else {
			dialer = &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			}
		}

		conn, err := dialer.DialContext(
			context.TODO(),
			network,
			kafkaSetting.BootstrapServers,
		)

		if err != nil {
			log.Printf("Error kafka connection: %v", err)
			return err
		}

		kafkaConnections[kafkaSetting.Id] = conn
	}

	return nil
}

func GetKafkaConnection(kafkaConnectionId string) (*kafka.Conn, error) {
	if len(kafkaConnectionId) == 0 ||
		len(kafkaConnections) == 0 ||
		kafkaConnections[kafkaConnectionId] == nil {
		return nil, errors.New("Without connection")
	}

	return kafkaConnections[kafkaConnectionId], nil
}

func getKafkaNetworking(kafkaConnectionType connection_settings.KafkaConnectionType) string {
	switch kafkaConnectionType {
	case connection_settings.KafkaConnectionPlaintext:
		return "tcp"
	case connection_settings.KafkaConnectionSsl:
		return "tcp"
	default:
		return "tcp"
	}
}
