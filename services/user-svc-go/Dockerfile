FROM golang:1.23-alpine AS build
WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod ./
COPY . .
RUN go mod tidy

RUN go build -o /out/user

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/user /app/user
EXPOSE 3000
ENTRYPOINT ["/app/user"]
