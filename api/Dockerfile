FROM golang:1.22-alpine as build

WORKDIR /app

COPY . .
RUN go build -o cp-voting-tool
RUN go vet ./...

FROM alpine:3.19 as runtime

WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /app .

ENTRYPOINT [ "/app/cp-voting-tool" ]



