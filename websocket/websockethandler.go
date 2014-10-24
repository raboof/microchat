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
				log.Println("Got message", msg, user_repo.FetchUser("1"))
				chat.Publish(msg)
				continue
			}
			break
		}
		close(closedSession)
		log.Println("sockjs session closed")
	}

}
