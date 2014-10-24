package main

import (
	"fmt"
	"github.com/raboof/microchat/events"
	"github.com/raboof/microchat/userrepo"
	"github.com/raboof/microchat/websocket"
	"log"
	"net/http"
	"strings"
)

func handleUser(user_repo *userrepo.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.URL.Query().Get("sessionId")
		fmt.Fprintf(w, "{ \"username\": \""+user_repo.FetchUser(sessionId).Name+"\" }")
	}
}

func handleUsers(user_repo *userrepo.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := user_repo.FetchUsers()
		var total = make([]string, 0)
		for i := 0; i < len(users); i++ {
			total = append(total, "\""+users[i].Name+"\"")
		}
		fmt.Fprintf(w, "["+strings.Join(total, ", ")+"]")
	}
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ messages }")
}

func main() {
	/* pre-provision */
	user_repo := userrepo.NewUserRepo()
	user_repo.StoreUser(userrepo.NewUser("987654321", "Hans"))

	/* start listening for domain events in background */
	eventListener := events.NewDomainEventListener(user_repo)
	//eventListener.Start("10.0.0.157:9092")
	eventListener.Start("169.254.101.81:9092")

	/* start listening for web-events */
	http.HandleFunc("/api/user", handleUser(user_repo))
	http.HandleFunc("/api/users", handleUsers(user_repo))
	http.HandleFunc("/api/messages", handleMessages)
	http.Handle("/ws/", websocket.WebsocketHandler(user_repo))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
