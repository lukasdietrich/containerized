FROM python:3-alpine

	WORKDIR /app
	COPY VERSION .

	RUN apk --no-cache add --virtual build-deps build-base \
		&& CFLAGS=-s pip3 install fava==$(cat VERSION) \
		&& apk del build-deps

	RUN adduser -D -H -u 1234 fava
	
	EXPOSE 5000/tcp
	VOLUME /data
	USER fava

	ENV LC_ALL "C.UTF-8"
	ENV LANG "C.UTF-8"
	ENV FAVA_HOST "0.0.0.0"

	CMD [ "fava", "/data/ledger.beancount" ]
