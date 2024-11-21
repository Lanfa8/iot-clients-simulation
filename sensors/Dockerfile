FROM        golang:latest

ADD         index.go /go/src/index.go
ADD         go.mod /go/src/go.mod
ADD         go.sum /go/src/go.sum
CMD         cd /go/src && go mod download && go run index.go