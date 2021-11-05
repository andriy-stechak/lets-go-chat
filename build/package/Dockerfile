FROM golang:1.17
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build ./cmd/chat/entry.go

FROM golang:1.17
WORKDIR /go/bin
COPY --from=0 /go/src/app/entry ./
ENTRYPOINT ["./entry"]
