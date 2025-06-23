package connections

import (
	"crypto/tls"
	"errors"
	"time"
	"wrench/app/manifest/connection_settings"

	"github.com/segmentio/kafka-go"
)

type KafkaConnection struct {
	Id      string
	Brokers []string
	Dialer  *kafka.Dialer
}

var kafkaConnections map[string]*KafkaConnection
var kafkaWriter map[string]*kafka.Writer = make(map[string]*kafka.Writer)

func loadConnectionsKafka(kafkaSettings []*connection_settings.KafkaConnectionSettings) error {

	if len(kafkaSettings) > 0 && kafkaConnections == nil {
		kafkaConnections = make(map[string]*KafkaConnection)
	}

	for _, setting := range kafkaSettings {
		var dialer *kafka.Dialer

		if setting.ConnectionType == connection_settings.KafkaConnectionSsl {
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

		kafkaConnections[setting.Id] = &KafkaConnection{
			Id:      setting.Id,
			Brokers: []string{setting.BootstrapServers},
			Dialer:  dialer,
		}
	}

	return nil
}

func GetKafkaConnection(kafkaConnectionId string) (*KafkaConnection, error) {
	if len(kafkaConnections) == 0 ||
		len(kafkaConnections) == 0 ||
		kafkaConnections[kafkaConnectionId] == nil {
		return nil, errors.New("Without connection")
	}

	return kafkaConnections[kafkaConnectionId], nil
}

func GetKafkaWrite(kafkaConnectionId string, topicName string) (*kafka.Writer, error) {

	kafkaWriterKey := kafkaConnectionId + "__" + topicName

	writer := kafkaWriter[kafkaWriterKey]
	if writer != nil {
		return writer, nil
	}

	conn, err := GetKafkaConnection(kafkaConnectionId)

	if err != nil {
		return nil, err
	}

	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  conn.Brokers,
		Topic:    topicName,
		Balancer: &kafka.LeastBytes{},
		Dialer:   conn.Dialer,
	})
	kafkaWriter[kafkaWriterKey] = writer
	return writer, nil
}
