FROM golang:1.22.3-alpine as build

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY ./cmd/web/*.go ./

RUN go build -o gorfc

FROM alpine:3.9 AS build-release-stage

WORKDIR /

COPY --from=build /app/gorfc /app/gorfc

EXPOSE 4000

CMD [ "/app/gorfc" ]