Running

Setup dependencies:
    setup_deps.sh

Build:
    make install
    
Run server:
    $GOPATH/bin/microchat

.. will listen for web events on port :8088
/api/user: return all logged in users
/api/user?sessionId=<valid-session-id>: return name of currently logged-in user
/api/message?sessionId=<valid-session-id>: return messages of logged-in user

.. will connect to kafka event-queue on 169.254.101.81:9092 on topic "my_topic"
    UserLoggedIn: received from authentication-subsystem via kafka
	UserLoggedOut: received from authentication-subsystem via kafka
