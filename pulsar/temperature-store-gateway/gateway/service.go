package gateway

import "github.com/apache/pulsar-client-go/pulsar"

type Service struct {
	Host        string
	Port        int
	PulsarURL   string
	PulsarTopic string

	pulsarClient   pulsar.Client
	pulsarProducer pulsar.Producer
}
