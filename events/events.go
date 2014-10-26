// Events packages
//
// Reads events from kafka queue
// based on these events the internal model is composed
package events

import (
	"github.com/Shopify/sarama"
	"github.com/raboof/microchat/userrepo"
	"log"
)

type EventListenerI interface {
	ListenAndServe(addressPort string) error
}

type EventHandlerFunc func(key string, value string, topic string, partition int32, offset int64)

type KafkaEventListener struct {
	repo        userrepo.UserRepoI
	handleEvent EventHandlerFunc
}

func NewKafkaEventListener(handler EventHandlerFunc) *KafkaEventListener {
	listnr := new(KafkaEventListener)
	listnr.handleEvent = handler

	return listnr
}

func (listener *KafkaEventListener) listenForEvents(hostnamePort string, topic string, clientId string, consumerGroup string) error {

	/* create broker */
	mb := sarama.NewBroker(hostnamePort)
	mdr := new(sarama.MetadataResponse)
	mdr.AddBroker(mb.Addr(), mb.ID())
	mdr.AddTopicPartition(topic, 0, 2)

	/* create client */
	client, err := sarama.NewClient(clientId, []string{mb.Addr()}, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	/* create consumer */
	consumer, err := sarama.NewConsumer(client, topic, 0, consumerGroup, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer consumer.Close()
	defer mb.Close()

	log.Printf("Start listening for kafka events on %s", hostnamePort)
	for {
		event := <-consumer.Events()
		listener.handleEvent(string(event.Key), string(event.Value), event.Topic, event.Partition, event.Offset)
	}

	return nil
}

func (listener *KafkaEventListener) ListenAndServe(hostnamePort string) error {
	return listener.listenForEvents(hostnamePort, "my_topic", "client_id", "my_consumer_group")
}
