# meshcentral

## docker-compose.yml

```
version: '3'

services:
	meshcentral:
		container_name: 'meshcentral'
		image: 'ghcr.io/lukasdietrich/dockerized-meshcentral'
		restart: 'always'
		logging:
			driver: 'journald'

		volumes:
			- './volumes/meshcentral-data:/app/meshcentral-data'
			- './volumes/meshcentral-files:/app/meshcentral-files'
```
