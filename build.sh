GO111MODULE=on  # Enable module mode
PATH="$PATH:$(go env GOPATH)/bin"
protoc -I proto/ proto/passman_api.proto --go_out=plugins=grpc:proto
