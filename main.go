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

func handleUser(user_repo *userrepo.UserRepo) http.HandlerFunc {
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
			user := user_repo.FetchUser(sessionId)
			if user == nil {
				http.Error(w, http.StatusText(404), 404)
			} else {
				fmt.Fprintf(w, "{ \"username\": \""+user.Name+"\" }")
			}
		}
	}
}

func handleMessage(user_repo *userrepo.UserRepo, forwarder *forwarder.Forwarder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method:%s, url:%s", r.Method, "messages")
		if r.Method == "GET" {
			sessionId := r.URL.Query().Get("sessionId")
			if sessionId != "" {
				user := user_repo.FetchUser(sessionId)
				if user == nil {
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
					user := user_repo.FetchUser(sessionId)
					if user == nil {
						http.Error(w, http.StatusText(404), 404)
					} else {
						log.Printf("Forwarding msg %s\n", messageText)
						msg := userrepo.NewMessage(sessionId, messageText)
						forwarder.Forward(msg)
					}
				}
			}
		}

	}
}

func main() {
	/* cenral store of users and their messages */
	var user_repo *userrepo.UserRepo
	user_repo = userrepo.NewUserRepo()

	/* pre-provision store for easy testing */
	user_repo.StoreUser(userrepo.NewUser("5678", "Hans"))
	user_repo.StoreUser(userrepo.NewUser("1234", "Grietje"))

	/* start listening for domain events in background */
	eventListener := events.NewDomainEventListener(user_repo)
	//eventListener.Start("10.0.0.157:9092")
	eventListener.Start("169.254.101.81:9092")

	/* forwarder  responssible for forwarding messages to logged in users */
	forwarder := forwarder.NewForwarder(user_repo)

	/* start listening for web-events */
	http.HandleFunc("/api/user", handleUser(user_repo))
	http.HandleFunc("/api/message", handleMessage(user_repo, forwarder))
	http.Handle("/ws/", websocket.WebsocketHandler(user_repo))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Start listening for web events at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
