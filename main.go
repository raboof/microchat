package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/raboof/microchat/userrepo"
	"github.com/raboof/microchat/websocket"
)

func handleUser(user_repo *userrepo.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.URL.Query().Get("sessionId")
		fmt.Fprintf(w, "{ \"username\": \""+user_repo.FetchUser(sessionId).Name+"\" }");
	}
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ users }")
}

func handleMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ messages }")
}


func main() {
	user_repo := userrepo.NewUserRepo()
	user_repo.StoreUser(userrepo.NewUser("1", "name 1"))
	user_repo.StoreUser(userrepo.NewUser("2", "name 2"))
	user_repo.StoreUser(userrepo.NewUser("3", "name 3"))

	http.HandleFunc("/api/user", handleUser(user_repo))
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/messages", handleMessages)
	http.Handle("/api/ws", websocket.WebsocketHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
