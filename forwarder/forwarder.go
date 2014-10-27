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
	log.Printf("ForwardMsgSent %s", msg.MessageText)
}

func (this *Forwarder) ForwardUserLoggedIn(user userrepo.User) {
	log.Printf("ForwardUserLoggedIn %s", user.Name)
}

func (this *Forwarder) ForwardUserLoggedOut(user userrepo.User) {
	log.Printf("ForwardUserLoggedOut %s", user.Name)
}
