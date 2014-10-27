package main

// Command-line application implementing chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/raboof/microchat/events"
	"github.com/raboof/microchat/forwarder"
	"github.com/raboof/microchat/userrepo"
	"github.com/raboof/microchat/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

func handleUser(user_repo userrepo.UserRepoI, forwarder forwarder.ForwarderI) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			sessionId := r.URL.Query().Get("sessionId")
			if strings.TrimSpace(sessionId) != "" {
				user, exists := user_repo.FetchUser(sessionId)
				if exists == false {
					http.Error(w, http.StatusText(404), 404)
				} else {
					jsonData, err := json.Marshal(user)
					if err != nil {
						http.Error(w, http.StatusText(500), 500)
					}
					fmt.Fprintf(w, string(jsonData))
				}
			} else {
				users := user_repo.FetchUsers()
				userNames := make([]string, 0, len(users))
				for i := range users {
					userNames = append(userNames, users[i].Name)
				}
				jsonData, err := json.Marshal(userNames)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				fmt.Fprintf(w, string(jsonData))
			}
		} else if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, http.StatusText(400), 400)
			} else {
				sessionId := r.PostForm.Get("sessionId")
				messageText := r.PostForm.Get("messageText")
				if strings.TrimSpace(sessionId) == "" || strings.TrimSpace(messageText) == "" {
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

func startEventListener(user_repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	eventListener := events.NewKafkaEventListener(handleEvent(user_repo, frwrdr))
	eventListener.ConnectAndReceive("169.254.101.81:9092")
}

func startWebServer(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	// start listening for web-events
	http.HandleFunc("/api/user", handleUser(repo, frwrdr))
	http.Handle("/ws/", websocket.WebsocketHandler(repo))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Start listening for web events at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

func readComamndLineInput(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	bio := bufio.NewReader(os.Stdin)

	var idx int64
	for idx = 0; ; idx++ {
		fmt.Printf("cli> ")
		line, _, err := bio.ReadLine()
		if err != nil {
			break
		}
		parts := strings.Split(string(line), " ")
		if len(parts) >= 1 && len(strings.TrimSpace(parts[0])) > 0 {
			cmd := fmt.Sprintf("%s,user_%d,uid", parts[0], idx)
			if parts[0] == "UserLoggedIn" {
				handleEvent(repo, frwrdr)("", cmd, "user", 1, idx)
			} else if parts[0] == "UserLoggedOut" {
				handleEvent(repo, frwrdr)("", cmd, "user", 1, idx)
			} else if parts[0] == "ChatMessage" {
				frwrdr.ForwardMsgSent(*userrepo.NewMessage("uid",
					fmt.Sprintf("test message %d", idx)))
			} else {
				fmt.Printf("Commands:\n\t%s : %s\n\t%s : %s\n\t%s : %s\n",
					"UserLoggedIn", "simulate logged-in event",
					"ChatMessage", "simulate chat-msg",
					"UserLoggedOut", "simulatelogged-out event")
			}
		}
	}
}

func main() {
	// cenral store of users and their messages
	// pre-provision store for easy testing
	userRepo := userrepo.NewUserRepo()
	userRepo.StoreUser(userrepo.NewUser("5678", "Hans"))
	userRepo.StoreUser(userrepo.NewUser("1234", "Grietje"))

	// forwarder is responssible for forwarding messages to other parts of the application
	forwarder := forwarder.NewForwarder(userRepo)

	// start listening for domain events in background
	go startEventListener(userRepo, forwarder)

	// start serving web requests
	go startWebServer(userRepo, forwarder)

	// read commands on stdin
	readComamndLineInput(userRepo, forwarder)

}
