package main

// Command-line application implementing chat

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/raboof/microchat/events"
	"github.com/raboof/microchat/forwarder"
	"github.com/raboof/microchat/userrepo"
	"github.com/raboof/microchat/websocket"
	"net/http"

	"log"
	"os"
	"strings"
)

func processChatMessage(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI,
	sender userrepo.User, msgText string) {
	msg := userrepo.NewMessage(sender.SessionId, msgText)
	sender.AddMsgSent(msg)
	repo.StoreUser(&sender)
	users := repo.FetchUsers()
	for _, rcver := range users {
		if rcver.SessionId != sender.SessionId {
			// store for fetching from UI
			rcver.AddMsgReceived(msg)
			repo.StoreUser(&rcver)

		}
	}
	//forward to other interested parties
	frwrdr.ForwardMsgSent(*msg)
}

func onPostMessage(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI, sessionId string, messageText string) (int, string) {
	if strings.TrimSpace(sessionId) == "" || strings.TrimSpace(messageText) == "" {
		return 400, "Invalid parameters"
	} else {
		user, exists := repo.FetchUser(sessionId)
		if exists == false {
			return 404, "Not found"
		} else {
			processChatMessage(repo, frwrdr, user, messageText)
			return 200, ""
		}
	}
}

func onGetUser(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI, sessionId string) (int, string) {
	if strings.TrimSpace(sessionId) == "" {
		return 400, "Invalid parameters"
	} else {
		user, exists := repo.FetchUser(sessionId)
		if exists == false {
			return 404, "Not found"
		} else {
			jsonData, err := json.Marshal(user)
			if err != nil {
				return 00, "Marshalling error"
			}
			return 200, string(jsonData)
		}
	}
}

func onGetAllUsers(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) (int, string) {
	users := repo.FetchUsers()
	userNames := make([]string, 0, len(users))
	for i := range users {
		userNames = append(userNames, users[i].Name)
	}
	jsonData, err := json.Marshal(userNames)
	if err != nil {
		return 500, "Marshalling error"
	}
	return 200, string(jsonData)
}

func listenAndServerRest(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	router := gin.Default()
	apiv2 := router.Group("/apiv2")
	{
		apiv2.POST("/usersession/:sessionId/message", func(c *gin.Context) {
			sessionId := c.Params.ByName("sessionId")
			messageText := c.Params.ByName("messageText")
			c.String(onPostMessage(repo, frwrdr, sessionId, messageText))
		})
		apiv2.GET("/usersession/:sessionId", func(c *gin.Context) {
			sessionId := c.Params.ByName("sessionId")
			c.String(onGetUser(repo, frwrdr, sessionId))
		})
		apiv2.GET("/usersession", func(c *gin.Context) {
			c.String(onGetAllUsers(repo, frwrdr))
		})
	}
	log.Println("Start listening for rest events at localhost:8089...")
	router.Run(":8089")
}

func onEvent(user_repo userrepo.UserRepoI, forwarder forwarder.ForwarderI) events.EventHandlerFunc {

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
	eventListener := events.NewKafkaEventListener(onEvent(user_repo, frwrdr))
	eventListener.ConnectAndReceive("169.254.101.81:9092")
}

func readCommand(r *bufio.Reader) (string, string, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return "", "", err
	}
	parts := strings.Split(string(line), " ")
	if len(parts) < 1 || len(strings.TrimSpace(parts[0])) == 0 {
		return "", "", errors.New("incomplete command")
	}
	if len(parts) >= 2 && len(strings.TrimSpace(parts[1])) > 0 {
		return parts[0], parts[1], nil
	}
	return parts[0], "", nil
}

func startCommandLineReader(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	reader := bufio.NewReader(os.Stdin)

	var idx int64
	for idx = 0; ; idx++ {
		fmt.Printf("cli> ")
		command, uid, err := readCommand(reader)
		if err == nil {
			if len(strings.TrimSpace(uid)) == 0 {
				uid = "ABC"
			}
			cmd := fmt.Sprintf("%s,%s,%s", command, uid, uid)
			user, userExists := repo.FetchUser(uid)
			if command == "UserLoggedIn" {
				if userExists == true {
					fmt.Printf("User %s already logged in\n", uid)
				} else {
					onEvent(repo, frwrdr)("", cmd, "user", 1, idx)
				}
			} else if command == "UserLoggedOut" {
				if userExists == false {
					fmt.Printf("User %s not logged in\n", uid)
				} else {
					onEvent(repo, frwrdr)("", cmd, "user", 1, idx)
				}
			} else if command == "ChatMessage" {
				if userExists == false {
					fmt.Printf("User %s not logged in\n", uid)
				} else {
					processChatMessage(repo, frwrdr, user, fmt.Sprintf("test message %d", idx))
				}
			} else {
				fmt.Printf("Commands:\n\t%s : %s\n\t%s : %s\n\t%s : %s\n",
					"UserLoggedIn [sessionName]", "simulate logged-in event",
					"ChatMessage [sessionName]", "simulate chat-msg",
					"UserLoggedOut [sessionName]", "simulatelogged-out event")
			}
		}
	}
}

func listenAndServerWeb(repo userrepo.UserRepoI, frwrdr forwarder.ForwarderI) {
	// start listening for web-events
	http.Handle("/ws/", websocket.WebsocketHandler(repo))
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Start listening for web events at localhost:8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
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

	// start listening using gin
	go listenAndServerWeb(userRepo, forwarder)
	go listenAndServerRest(userRepo, forwarder)

	// read commands on stdin
	startCommandLineReader(userRepo, forwarder)
}
