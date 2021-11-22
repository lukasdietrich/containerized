FROM golang:alpine as build

	WORKDIR /github.com/lukasdietrich/dockerized/updater
	COPY . .

	RUN go build ./cmd/dockerized-updater

FROM alpine:latest

	WORKDIR /app
	COPY --from=build /github.com/lukasdietrich/dockerized/updater/dockerized-updater /app/

	RUN apk update \
		&& apk upgrade \
		&& apk --no-cache add ca-certificates

	ENV DOCKERIZED_SCHEDULE=TODO
	ENV DOCKERIZED_DEPLOY_KEY=TODO
	ENV DOCKERIZED_GIT_AUTHOR_NAME=TODO
	ENV DOCKERIZED_GIT_AUTHOR_EMAIL=TODO

	CMD [ "/app/dockerized-updater" ]

