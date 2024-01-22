FROM golang:1.21.6-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o bin/parser ./main.go

# add the binary to the path
ENV PATH=$PATH:/app/bin

CMD ["./bin/parser", "stats"]