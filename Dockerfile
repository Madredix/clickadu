FROM golang:1.13 AS build-env
# Add project directory to Docker image.
ADD . /go/src/github.com/Madredix/clickadu
WORKDIR /go/src/github.com/Madredix/clickadu

RUN apt-get update && apt-get upgrade
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger
ENV GOOS=linux
ENV GARCH=amd64
ENV CGO_ENABLED=0
RUN make build
EXPOSE 2000
CMD ./app


# final stage
FROM alpine
WORKDIR /app
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build-env /go/src/github.com/Madredix/clickadu/ /app/
COPY --from=build-env /go/src/github.com/Madredix/clickadu/config.json.example config.json
CMD ./app
