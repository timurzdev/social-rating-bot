FROM golang:latest
RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get install -y postgresql-client

RUN go mod download
RUN go build -o social-rating-bot ./cmd/main.go

CMD ["./social-rating-bot"]
