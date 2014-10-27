package main

// Command-line application implementing chat

import (
	"fmt"
	"github.com/raboof/microchat/events"
	"github.com/raboof/microchat/forwarder"
	"github.com/raboof/microchat/userrepo"
	"github.com/raboof/microchat/websocket"
	"log"
	"net/http"
	"strings"
)

func handleUser(user_repo userrepo.UserRepoI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method:%s, url:%s", r.Method, "users")
		sessionId := r.URL.Query().Get("sessionId")
		if sessionId == "" {
			/* fetch all users */
			users := user_repo.FetchUsers()
			var total = make([]string, 0)
			for i := 0; i < len(users); i++ {
				total = append(total, "\""+users[i].Name+"\"")
			}
			fmt.Fprintf(w, "["+strings.Join(total, ", ")+"]")
		} else {
			/* fetch single user */
			user, exists := user_repo.FetchUser(sessionId)
			if exists == false {
				http.Error(w, http.StatusText(404), 404)
			} else {
				fmt.Fprintf(w, "{ \"username\": \""+user.Name+"\" }")
			}
		}
	}
}

func handleMessage(user_repo userrepo.UserRepoI, forwarder *forwarder.Forwarder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method:%s, url:%s", r.Method, "messages")
		if r.Method == "GET" {
			sessionId := r.URL.Query().Get("sessionId")
			if sessionId != "" {
				user, exists := user_repo.FetchUser(sessionId)
				if exists == false {
					http.Error(w, http.StatusText(404), 404)
				} else {
					fmt.Fprintf(w, "{ \"name\":\"%s\", \"receivedMsgCount\": %d, \"sentMsgCount\":%d }",
						user.Name,
						len(user.ReceivedMessages),
						len(user.SentMessages))
				}
			} else {
				http.Error(w, http.StatusText(404), 404)
			}
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err == nil {
				sessionId := r.PostForm.Get("sessionId")
				messageText := r.PostForm.Get("messageText")
				if sessionId == "" || messageText == "" {
					http.Error(w, http.StatusText(400), 400)
				} else {
					_, exists := user_repo.FetchUser(sessionId)
					if exists == false {
						http.Error(w, http.StatusText(404), 404)
					} else {
						log.Printf("Forwarding msg %s\n", messageText)
						msg := userrepo.NewMessage(sessionId, messageText)
						forwarder.ForwardMsgSent(*msg)
					}
				}
			}
		}

	}
}

func handleEvent(user_repo userrepo.UserRepoI, forwarder forwarder.ForwarderI) events.EventHandlerFunc {

	return func(key string, value string, topic string, partition int32, offset int64) {

		log.Printf("Received cosumer event with key:'%s', value:'%s', topic:'%s', partition: %d, offset: %d",
			key, value, topic, partition, offset)

		s := strings.Split(string(value), ",")
		if len(s) < 3 {
			log.Printf("Event incomplete: '%s'", value)
		} else {
			eventName, userName, sessionId := s[0], s[1], s[2]
			user := userrepo.NewUser(sessionId, userName)
			if eventName == "UserLoggedIn" {
				user_repo.StoreUser(user)
				forwarder.ForwardUserLoggedIn(*user)
			} else if eventName == "UserLoggedOut" {
				user_repo.RemoveUser(user)
				forwarder.ForwardUserLoggedOut(*user)
			} else {
				log.Printf("Unrecognized event %s", eventName)
			}
		}
	}
}

func main() {
	// cenral store of users and their messages
	user_repo := userrepo.NewUserRepo()

	// pre-provision store for easy testing
	user_repo.StoreUser(userrepo.NewUser("5678", "Hans"))
	user_repo.StoreUser(userrepo.NewUser("1234", "Grietje"))

	// forwarder is responssible for forwarding messages to other parts of the application
	forwarder := forwarder.NewForwarder(user_repo)

	// start listening for domain events in background
	eventListener := events.NewKafkaEventListener(handleEvent(user_repo, forwarder))
	go eventListener.ConnectAndReceive("169.254.101.81:9092")

	// start listening for web-events
	http.HandleFunc("/api/user", handleUser(user_repo))
	http.HandleFunc("/api/message", handleMessage(user_repo, forwarder))
	http.Handle("/ws/", websocket.WebsocketHandler(user_repo))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Start listening for web events at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
