FROM golang:1.22.3

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY *.go ./

RUN go build -o gorfc

CMD [ "/app/gorfc" ]