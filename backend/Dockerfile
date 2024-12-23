FROM golang:latest

RUN mkdir /base
WORKDIR /base

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

CMD ./run.sh