package events

import (
    "log"
    "github.com/Shopify/sarama"
    "github.com/raboof/microchat/userrepo"
)

type DomainEventListener struct {
     quit chan bool
     repo *userrepo.UserRepo
}


func NewDomainEventListener( repo *userrepo.UserRepo) *DomainEventListener {
    listnr := new(DomainEventListener)
    listnr.repo = repo

    return listnr
}

func (listener *DomainEventListener) listenForEvents( hostnamePort string, topic string, clientId string, consumerGroup string) error {
	mb := sarama.NewBroker(hostnamePort)

        log.Printf( "Starting domain event listener on %s", hostnamePort )

	mdr := new(sarama.MetadataResponse)
	mdr.AddBroker(mb.Addr(), mb.ID())
	mdr.AddTopicPartition(topic, 0, 2)

	for i := 0; i < 10; i++ {
		fr := new(sarama.FetchResponse)
		fr.AddMessage(topic, 0, nil, sarama.ByteEncoder([]byte{0x00, 0x0E}), int64(i))
	}

	client, err := sarama.NewClient(clientId, []string{mb.Addr()}, nil)

	if err != nil {
		return err
	}

	defer client.Close()

	consumer, err := sarama.NewConsumer(client, topic, 0, consumerGroup, nil)
	if err != nil {
		return err
	}
	defer consumer.Close()
	defer mb.Close()

	event := <-consumer.Events()

	log.Println("ja")
	log.Println(event.Offset)

	return nil
}

func (listener *DomainEventListener) Start( hostnamePort string ) {
    go listener.listenForEvents( hostnamePort, "my_topic","client_id", "my_consumer_group" )
}

func (listener *DomainEventListener) Stop() {
}

