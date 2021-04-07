protofiles:
	@find internal/pb -type f -name '*.go' -delete
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/pb/common.proto internal/pb/room/room.proto internal/pb/message/message.proto

zapcheck:
	@command -v zapw > /dev/null 2>&1 || GO111MODULE=off go get github.com/sethvargo/zapw/cmd/zapw
	@zapw ./...