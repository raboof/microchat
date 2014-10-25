package forwarder

import (
	"github.com/raboof/microchat/userrepo"
	"log"
)

type ForwarderI interface {
	Forward(msg *userrepo.Message)
}

type Forwarder struct {
	repo userrepo.UserRepoI
}

func NewForwarder(repo userrepo.UserRepoI) *Forwarder {
	frwrdr := new(Forwarder)
	frwrdr.repo = repo

	return frwrdr
}

func (frwrdr *Forwarder) Forward(msg *userrepo.Message) {
	sender := frwrdr.repo.FetchUser(msg.OriginatorSessionId)
	if sender != nil {
		log.Printf("Adding msg to sender %s", sender.Name)
		sender.AddMsgSent(msg)
		users := frwrdr.repo.FetchUsers()
		for _, rcver := range users {
			if rcver.SessionId != sender.SessionId {
				log.Printf("Adding msg to receiver %s", rcver.Name)
				/* store for fetching from UI */
				rcver.AddMsgReceived(msg)

				/* TODO: forward to web-socket */
			}
		}
	}
}
