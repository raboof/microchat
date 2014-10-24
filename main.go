package main

import (
	"log"
	"fmt"
	"net/http"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ \"username\": \"raboof\" }");
}
func handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ users }");
}
func handleMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{ messages }");
}

func main() {
	http.HandleFunc("/api/user", handleUser)
	http.HandleFunc("/api/users", handleUsers)
	http.HandleFunc("/api/messages", handleMessages)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:8080...")
	http.ListenAndServe(":8080", nil)
}
