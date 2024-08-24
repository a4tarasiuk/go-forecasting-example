FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN go build -o /go-docker

EXPOSE 8080

CMD [ "/go-docker" ]
