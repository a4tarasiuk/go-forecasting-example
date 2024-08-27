FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN go build -C cmd/ -o /go-forecasting

EXPOSE 8080

CMD [ "/go-forecasting" ]
