package websocket

import (
	"github.com/igm/pubsub"
	"github.com/raboof/microchat/userrepo"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

var chat pubsub.Publisher

func WebsocketHandler(user_repo userrepo.UserRepoI) http.Handler {
	return sockjs.NewHandler("/ws", sockjs.DefaultOptions, echoHandler(user_repo))
}

func echoHandler(user_repo userrepo.UserRepoI) func(sockjs.Session) {
	users := make(map[string]userrepo.User)
	return func(session sockjs.Session) {
		log.Println("new sockjs session established")
		var closedSession = make(chan struct{})
		defer func() {
			user, ok := users[session.ID()]
			if ok == true {
				chat.Publish("[info] " + user.Name + " left chat")
			}
		}()
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
				user, ok := users[session.ID()]
				if ok == false {
					user, exists := user_repo.FetchUser(msg)
					if exists == false {
						log.Println("Not a user id", msg)
						break
					}
					chat.Publish("[info] " + user.Name + " joined chat")
					users[session.ID()] = user
				} else {
					chat.Publish(user.Name + ": " + msg)
				}
				continue
			}
			break
		}
		close(closedSession)
		log.Println("sockjs session closed")
	}

}
