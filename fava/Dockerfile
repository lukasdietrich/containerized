FROM python:3-alpine

	WORKDIR /app
	COPY VERSION .
	RUN apk --no-cache add --virtual build-deps build-base \
		&& pip3 install fava==$(cat VERSION) \
		&& apk del build-deps
	
	EXPOSE 5000/tcp
	VOLUME "/data"
	USER nobody
	ENV "FAVA_HOST"="0.0.0.0"

	CMD [ "fava", "/data/ledger.beancount" ]
