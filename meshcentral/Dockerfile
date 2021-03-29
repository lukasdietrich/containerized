FROM node:alpine

	WORKDIR /app
	RUN npm install meshcentral

	EXPOSE 80/tcp 443/tcp

	VOLUME [ "/app/meshcentral-data", "/app/meshcentral-files" ]

	CMD [ "node", "node_modules/meshcentral" ]
