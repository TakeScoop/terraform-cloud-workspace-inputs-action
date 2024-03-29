FROM golang:1.18-alpine AS build

RUN apk add git

WORKDIR /src
COPY . ./

RUN go build

ENTRYPOINT ["/src/compose-inputs"]
