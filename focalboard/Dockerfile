FROM node:alpine as build-frontend

	WORKDIR /build
	COPY submodules/focalboard .

	RUN apk --no-cache add make \
		&& make prebuild webapp

FROM golang:alpine as build-backend
	
	WORKDIR /build
	COPY submodules/focalboard .

	RUN apk --no-cache add build-base git \
		&& make server

FROM alpine:latest

	WORKDIR /app
	COPY --from=build-frontend /build/webapp/pack /app/pack
	COPY --from=build-backend /build/bin/focalboard-server /app/focalboard
	COPY files /app
	
	EXPOSE 8000/tcp
	VOLUME [ "/data" ]

	CMD [ "/app/focalboard" ]
