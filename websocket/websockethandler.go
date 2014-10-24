package websocket

import (
	"log"
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"github.com/raboof/microchat/userrepo"
	"net/http"
)

var chat pubsub.Publisher

func WebsocketHandler(user_repo *userrepo.UserRepo) http.Handler {
	return sockjs.NewHandler("/ws", sockjs.DefaultOptions, echoHandler(user_repo))
}

func echoHandler(user_repo *userrepo.UserRepo) func(sockjs.Session) {
	users := make(map[string]*userrepo.User)
	return func(session sockjs.Session) {
		log.Println("new sockjs session established")
		var closedSession = make(chan struct {})
		chat.Publish("[info] new participant joined chat")
		defer chat.Publish("[info] participant left chat")
		go func() {
			reader, _ := chat.SubChannel(nil)
			for {
				select {
				case <-closedSession:
					return
				case msg := <-reader:
					if err := session.Send(msg.(string)); err != nil {
						return
					}
				}

			}
		}()
		for {
			if msg, err := session.Recv(); err == nil {
				user := users[session.ID()]
				if (user == nil) {
					user = user_repo.FetchUser(msg)
					if (user == nil) {
						log.Println("Not a user id", msg)
						break
					}
					users[session.ID()] = user
				} else {
					chat.Publish(user.Name + ":" + msg)
				}
				continue
			}
			break
		}
		close(closedSession)
		log.Println("sockjs session closed")
	}

}
