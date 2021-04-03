# meshcentral

<https://github.com/Ylianst/MeshCentral>

## Documentation

- <https://info.meshcentral.com/downloads/MeshCentral2/MeshCentral2InstallGuide.pdf>
- <https://info.meshcentral.com/downloads/MeshCentral2/MeshCentral2UserGuide.pdf>

## docker-compose.yml

```yaml
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
