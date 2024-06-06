FROM golang:1.22.3 as build

WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY *.go ./

RUN go build -o gorfc

# FROM gcr.io/distroless/base-debian11 AS build-release-stage

# WORKDIR /

# COPY --from=build /app/gorfc /app/gorfc

EXPOSE 4000

CMD [ "/app/gorfc" ]