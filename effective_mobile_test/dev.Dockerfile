FROM golang:1.22.2

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app

RUN mkdir "/build"

COPY . .

COPY . .

RUN go get github.com/githubnemo/CompileDaemon@v1.4.0
RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

ENTRYPOINT CompileDaemon -build="go build -o /build/app" -command="/build/app"