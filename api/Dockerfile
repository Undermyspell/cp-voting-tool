FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN apk add npm
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    # npm run build && \
    templ generate && \
    go generate

RUN go build -o cp-voting-tool
RUN go vet ./...

FROM alpine:3.20 AS runtime

WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /app .

ENTRYPOINT [ "/app/cp-voting-tool" ]



