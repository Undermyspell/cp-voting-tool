# Conplement Voting Tool
This is the repository of the conplement voting tool used in employee meetings

## Web Frontend 1 (built with svelte) ⚡️
- [Svelte](https://svelte.dev/)
- [Vite](https://vitejs.dev/)

## Web Frontend 2 (built with go, templ, htmx, alpinejs) ⚡️
- [Go](https://go.dev/)
- [Templ](https://templ.guide/)
- [htmx](https://htmx.org/)
- [Alpine.js](https://alpinejs.dev/)

## Api Backend (built with go) ⚡️
### Web api framework
- [Go](https://go.dev/)
- [Gin-Gonic](https://github.com/gin-gonic/gin)
### Testing
- [Testify](https://github.com/stretchr/testify)
- [Testcontainers](https://testcontainers.com/)
### Live updates
- [Centrifuge](https://github.com/centrifugal/centrifuge)
### Data storage
#### Redis (previously in use)
- [Redis-Cloud](https://app.redislabs.com/#/)
- [Go-ReJSON](https://github.com/nitishm/go-rejson)
#### PostgreSQL (currently in use)
- [Neon PostgreSQL](https://neon.tech/)


## Get Started

### Create an environment variables file

You can create a env.yaml file in the ./api directory by creating it manually from the env.template.yaml. This file contains no secrets. This application can use different storage types
- In Memory
- Redis
- PostgreSQL

You can configure the type in the env.yaml file, yet You have to fill in the PostgreSql connection string or redis endpoint

If you want to start the api with a redis or postgreSQL cloud storage you have to decrypt the env.enc.yaml using [Mozilla sops](https://github.com/getsops/sops). Please contact the developer team for further information.

### Run docker-compose file

If a env.yaml file is located in the ./api directory you can run the application using the docker-compose file located at the root of the project.

```shell
docker-compose up -d
```