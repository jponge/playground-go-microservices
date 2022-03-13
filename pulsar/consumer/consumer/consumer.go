package consumer

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
)

type Consumer struct {
	PulsarURL   string
	PulsarTopic string
}

func (consumer Consumer) Start() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: consumer.PulsarURL,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reader, err := client.CreateReader(pulsar.ReaderOptions{
		Topic:          consumer.PulsarTopic,
		StartMessageID: pulsar.EarliestMessageID(),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	fmt.Println("🚀")
	for {
		message, err := reader.Next(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v -- %s -- %s\n", message.ID(), message.EventTime(), string(message.Payload()))
	}
}
