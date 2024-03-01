# Conplement Voting Tool
This is the repository of the conplement voting tool used in employee meetings

## Web Frontend (built with svelte) ⚡️
- [Svelte](https://svelte.dev/)
- [Vite](https://vitejs.dev/)

## Api Backend (built with go) ⚡️
- [Go](https://go.dev/)
- [Gin-Gonic](https://github.com/gin-gonic/gin)
- [Redis-Cloud](https://app.redislabs.com/#/)
- [Go-ReJSON](https://github.com/nitishm/go-rejson)
- [Testify](https://github.com/stretchr/testify)
- [Centrifuge](https://github.com/centrifugal/centrifuge)

## Get Started

### Create an environment variables file

You can create a env.yaml file in the ./api directory by creating it manually from the env.template.yaml. This file contains no secrets and starts the api with an inMemory store. 

If you want to start the api with a redis cloud storage you have to decrypt the env.enc.yaml using [Mozilla sops](https://github.com/getsops/sops). Please contact the developer team for further information.

### Run docker-compose file

If a env.yaml file is located in the ./api directory you can run the application using the docker-compose file located at the root of the project.

```shell
docker-compose up -d
```