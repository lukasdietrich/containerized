FROM node:alpine

	WORKDIR /app
	COPY VERSION .
	RUN npm install meshcentral@$(cat VERSION) --loglevel verbose

	EXPOSE 80/tcp 443/tcp

	VOLUME [ "/app/meshcentral-data", "/app/meshcentral-files" ]

	CMD [ "node", "node_modules/meshcentral" ]
