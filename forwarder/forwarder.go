package forwarder

import (
	"github.com/raboof/microchat/userrepo"
	"log"
)

type ForwarderI interface {
	ForwardUserLoggedIn(user userrepo.User)
	ForwardUserLoggedOut(user userrepo.User)
	ForwardMsgSent(msg userrepo.Message)
}

type Forwarder struct {
	repo userrepo.UserRepoI
}

func NewForwarder(repo userrepo.UserRepoI) *Forwarder {
	frwrdr := new(Forwarder)
	frwrdr.repo = repo

	return frwrdr
}

func (this *Forwarder) ForwardMsgSent(msg userrepo.Message) {
	sender, exists := this.repo.FetchUser(msg.OriginatorSessionId)
	if exists == true {
		log.Printf("Adding msg to sender %s", sender.Name)
		sender.AddMsgSent(&msg)
		users := this.repo.FetchUsers()
		for _, rcver := range users {
			if rcver.SessionId != sender.SessionId {
				log.Printf("Adding msg to receiver %s", rcver.Name)
				/* store for fetching from UI */
				rcver.AddMsgReceived(&msg)

				/* TODO: forward to web-socket */
			}
		}
	}
}

func (this *Forwarder) ForwardUserLoggedIn(user userrepo.User) {
}

func (this *Forwarder) ForwardUserLoggedOut(user userrepo.User) {
}
