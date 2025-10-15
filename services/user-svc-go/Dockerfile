FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /out/user

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/user /app/user
EXPOSE 3000
ENTRYPOINT ["/app/user"]
