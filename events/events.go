package events

import (
	"github.com/Shopify/sarama"
	"github.com/raboof/microchat/userrepo"
	"log"
	"strings"
)

type DomainEventListener struct {
	quit chan bool
	repo *userrepo.UserRepo
}

func NewDomainEventListener(repo *userrepo.UserRepo) *DomainEventListener {
	listnr := new(DomainEventListener)
	listnr.repo = repo

	listnr.quit = make(chan bool)
	return listnr
}

func (listener *DomainEventListener) listenForEvents(hostnamePort string, topic string, clientId string, consumerGroup string) error {

	log.Printf("Starting domain event listener on %s", hostnamePort)

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

	consumer, err := sarama.NewConsumer(client, topic, 0, consumerGroup, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	defer consumer.Close()
	defer mb.Close()

	log.Printf("Listen for events")
	for {
		event := <-consumer.Events()
		listener.handleEvent(event)
	}

	return nil
}

func (listener *DomainEventListener) handleEvent(event *sarama.ConsumerEvent) {

	log.Printf("Received event offset: %d, topic: %s, value: '%s'", event.Offset, event.Topic, event.Value)

	s := strings.Split(string(event.Value), ",")
	if len(s) < 3 {
		log.Printf("Event has not enough parameters")
	} else {
		eventName, userName, sessionId := s[0], s[1], s[2]
		user := userrepo.NewUser(userName, sessionId)
		if eventName == "UserLoggedIn" {
			listener.repo.StoreUser(user)
		} else if eventName == "UserLoggedOut" {
			listener.repo.RemoveUser(user)
		} else {
			log.Printf("Unrecognized event %s", eventName)
		}
	}
}

func (listener *DomainEventListener) Start(hostnamePort string) {
	go listener.listenForEvents(hostnamePort, "my_topic", "client_id", "my_consumer_group")
}

func (listener *DomainEventListener) Stop() {
	listener.quit <- true
}
