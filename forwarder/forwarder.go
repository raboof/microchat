package forwarder

import (
	"github.com/raboof/microchat/userrepo"
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
        sender := frwrdr.repo.FetchUser( msg.OriginatorSessionId ) 
	if sender != nil {
		sender.AddMsgSent(msg)
		users := frwrdr.repo.FetchUsers()
		for _, user := range users {
			/* store for fetching from UI */
			user.AddMsgReceived(msg)

			/* TODO: forward to web-socket */
		}
	}
}
