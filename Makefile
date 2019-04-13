.PHONY: gomod build b-chat b-chat-with-mdns b-chat-with-rendezvous b-dht b-echo b-host b-http-proxy b-libp2p-host b-multipro b-raft b-relay b-routed-echo

default: build

gomod: export GO111MODULE=on
gomod:
	go mod tidy
	go mod vendor

build: gomod
build: b-chat b-chat-with-mdns b-chat-with-rendezvous b-dht b-echo b-host b-http-proxy b-libp2p-host b-multipro b-raft b-relay b-routed-echo

b-chat:
	CGO_ENABLED=0 go build -o bin/chat chat/chat.go

b-chat-with-mdns:
	CGO_ENABLED=0 go build -o bin/chat-with-mdns chat-with-mdns/*.go

b-chat-with-rendezvous:
	CGO_ENABLED=0 go build -o bin/chat-with-rendezvous chat-with-rendezvous/*.go

b-dht:
	CGO_ENABLED=0 go build -o bin/dht dht/cmd/*.go

b-echo:
	CGO_ENABLED=0 go build -o bin/echo echo/main.go

b-host:
	CGO_ENABLED=0 go build -o bin/host host/main.go

b-http-proxy:
	CGO_ENABLED=0 go build -o bin/http-proxy http-proxy/proxy.go

b-libp2p-host:
	CGO_ENABLED=0 go build -o bin/libp2p-host libp2p-host/host.go

b-multipro:
	CGO_ENABLED=0 go build -o bin/multipro multipro/*.go

b-raft:
	CGO_ENABLED=0 go build -o bin/raft raft/main.go

b-relay:
	CGO_ENABLED=0 go build -o bin/relay relay/main.go

b-routed-echo:
	CGO_ENABLED=0 go build -o bin/routed-echo routed-echo/*.go



