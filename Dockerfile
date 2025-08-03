FROM golang:alpine AS builder
WORKDIR /src/app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server .

FROM alpine
WORKDIR /root/
COPY --from=builder /src/app ./app
ENTRYPOINT ["./app/server"]