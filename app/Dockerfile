FROM golang:1.19.5-alpine

ADD . /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /simple-go-web

EXPOSE 8000

CMD ["/simple-go-web"]
