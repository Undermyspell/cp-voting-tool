FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN apk add npm && \
    npm install
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate && \
    go generate

RUN go build -o cp-voting-tool
RUN go vet ./...

FROM alpine:3.21 AS runtime

WORKDIR /app

ENV GIN_MODE=release

COPY --from=build /app .

CMD [ "/app/cp-voting-tool" ]



