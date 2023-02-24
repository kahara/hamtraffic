FROM golang:1.19-bullseye

WORKDIR /workdir
COPY . .
RUN go get .
ENTRYPOINT ["go", "run", "cmd/hamtraffic/main.go" ]
