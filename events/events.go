package events

import (
	"github.com/Shopify/sarama"
	"github.com/raboof/microchat/userrepo"
	"log"
	"strings"
)

type DomainEventListenerI interface {
	Start(addressPort string) error
	Stop()
	HandleUserLoggedIn(user *userrepo.User)
	HandleUserLoggedOut(user *userrepo.User)
}

type DomainEventListener struct {
	quit chan bool
	repo userrepo.UserRepoI
}

func NewDomainEventListener(repo userrepo.UserRepoI) *DomainEventListener {
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

	/* create consumer */
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
		log.Printf("Event incomplete: '%s'", event.Value)
	} else {
		eventName, userName, sessionId := s[0], s[1], s[2]
		user := userrepo.NewUser(sessionId, userName)
		if eventName == "UserLoggedIn" {
			listener.HandleUserLoggedIn(user)
		} else if eventName == "UserLoggedOut" {
			listener.HandleUserLoggedOut(user)
		} else {
			log.Printf("Unrecognized event %s", eventName)
		}
	}
}

func (listener *DomainEventListener) HandleUserLoggedIn(user *userrepo.User) {
	listener.repo.StoreUser(user)
}

func (listener *DomainEventListener) HandleUserLoggedOut(user *userrepo.User) {
	listener.repo.RemoveUser(user)
}

func (listener *DomainEventListener) Start(hostnamePort string) {
	go listener.listenForEvents(hostnamePort, "my_topic", "client_id", "my_consumer_group")
}

func (listener *DomainEventListener) Stop() {
	listener.quit <- true
}
