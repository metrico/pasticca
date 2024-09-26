FROM golang:1.23

WORKDIR /app
COPY ./go/go.mod ./
RUN go mod download

COPY ./go .
RUN go build -v -o /usr/local/bin/duckserver ./...

CMD ["duckserver"]
