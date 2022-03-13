package gateway

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/jponge/playground-go-microservices/pulsar/temperature-store-gateway/data"
	"log"
	"time"
)

func (service *Service) initPulsarClient() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               service.PulsarURL,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not create the Pulsar client: %v", err)
	}
	service.pulsarClient = client

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: service.PulsarTopic,
		//DisableBatching: true,
	})
	if err != nil {
		log.Fatalf("Could not create the Pulsar producer: %v", err)
	}
	service.pulsarProducer = producer

	log.Println("🚀 Pulsar client initialized")
}

func (service Service) closePulsarClient() {
	if service.pulsarProducer != nil {
		service.pulsarProducer.Close()
	}
	if service.pulsarClient != nil {
		service.pulsarClient.Close()
	}
}

func (service Service) pushToPulsar(ctx context.Context, payload data.Payload) error {
	jsonBytes, err := payload.ToJSON()
	if err != nil {
		return err
	}
	id, err := service.pulsarProducer.Send(ctx, &pulsar.ProducerMessage{
		Key:     payload.SensorID,
		Payload: jsonBytes,
	})
	if err != nil {
		return err
	}
	log.Printf("Sent to Pulsar with message id: %s", id)
	return nil
}
