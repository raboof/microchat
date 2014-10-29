## Introduction
This is a small POC to prove the microservice concept in the Go programming language. 
This "microchat"-application consists of two parts: 
 - authentication and 
 - chat. 

The idea is that unavailability of one component (example: authentication) does not teardown other components that depend on it (example: chat).

In order to achieve this, components do not communicate directly with each other. Instead they broadcast state-changes as "events" on a queue. Dependend components can consume these events and act on them.

### Authentication component
So in our casse authentication-component provides a UI for subscribing and logging in and out.
Successfull login and logout are broadcasted as "UserLoggedIn"  and "UserLoggedOut"-evens.
After successfull login, the customer will be directed to the chat application.
The chat application will build a list of active users with their messages in memory.
Nothing is persisted.

### Chat component
The chat-component provides a UI and restfull-service to be able to do chat. In addition this, it listen for user-events (=UserLoggedIn and UserLoggedOut events).

### Infrastructure
We use a Kafka-queue to distribute events from producer to consumers. 

## Running the application

### Setup dependencies:
    setup_deps.sh

### Test the application
    make test
    
### Install:
    make install

### Run server:
    $GOPATH/bin/microchat

## Interfaces

### Web interface
GET /

### Rest interface
    GET /apiv2/usersession/ : return all logged in users

    GET /apiv2/usersession/:sessionId : return currently logged-in user with its messages

    POST /apiv2/usersession/:sessionId/message : Send a new chat messages to all other participants of the chat

### Event interface

The application will connect to kafka event-queue on 169.254.101.81:9092 on topic "my_topic"
    
    UserLoggedIn: received from authentication-subsystem via kafka
    
    UserLoggedOut: received from authentication-subsystem via kafka
