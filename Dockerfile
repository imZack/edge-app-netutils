FROM golang:1.10-alpine as builder

WORKDIR /go/src/github.com/imZack/edge-app-netutils
COPY src ./
# RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go
RUN GOOS=linux GOARCH=arm GOARM=7 go build -o main main.go

FROM arm32v7/debian:9-slim

WORKDIR /
RUN apt-get update && apt-get install -y stress && rm -rf /var/lib/apt/lists/
COPY --from=builder /go/src/github.com/imZack/edge-app-netutils/main /usr/local/bin/main
RUN mkdir -p /var/www/html
COPY html /var/www/html

CMD [ "/usr/local/bin/main" ]