FROM golang:1.13 AS build-env
# Add project directory to Docker image.
ADD . /go/src/https://github.com/Madredix/clickadu
WORKDIR /go/src/https://github.com/Madredix/clickadu

RUN apk update && apk upgrade
RUN go get github.com/go-swagger/go-swagger
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
COPY --from=build-env /go/src/https://github.com/Madredix/clickadu/ /app/
COPY --from=build-env /go/src/https://github.com/Madredix/clickadu/config.json.exemple config.json
CMD ./app
