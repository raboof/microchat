#!/bin/sh -x

PKGS="github.com/Shopify/sarama \
      github.com/stretchr/testify/assert \
      github.com/gorilla/websocket/ \
      github.com/igm/pubsub \
      gopkg.in/igm/sockjs-go.v2/sockjs"

for pkg in ${PKGS}
do
   go get ${pkg}
   cd $GOPATH/src/${pkg}
   go clean; go fmt; go build; go install
done
