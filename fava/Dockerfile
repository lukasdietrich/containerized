FROM python:3-alpine

	WORKDIR /app
	COPY VERSION .

	RUN apk --no-cache add --virtual build-deps libxml2-dev libxslt-dev gcc musl-dev g++ bison \
		&& pip3 install fava==$(cat VERSION) \
		&& apk del build-deps

	RUN adduser -D -H -u 1234 fava
	
	EXPOSE 5000/tcp
	USER fava

	VOLUME /data
	VOLUME /plugins

	ENV LC_ALL "C.UTF-8"
	ENV LANG "C.UTF-8"
	ENV PYTHONPATH "${PYTHONPATH}:/plugins"
	ENV FAVA_HOST "0.0.0.0"

	CMD [ "fava", "/data/ledger.beancount" ]
