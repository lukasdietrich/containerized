FROM node:alpine

	WORKDIR /data
	RUN npm install -g meshcentral

	EXPOSE 80/tcp 443/tcp

	VOLUME [ "/data/meshcentral-data", "/data/meshcentral-files" ]

	CMD [ "meshcentral" ]
