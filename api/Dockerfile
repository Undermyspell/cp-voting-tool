FROM golang:1.21-alpine as build

WORKDIR /app

COPY . .
RUN go build -o cp-voting-tool
RUN go vet ./...

FROM alpine:3.18 as runtime

WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /app .

ENTRYPOINT [ "/app/cp-voting-tool" ]



