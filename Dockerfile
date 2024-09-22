FROM golang:1.22.2

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN chmod +x ./scripts/startup.sh
RUN chmod +x ./scripts/shutdown.sh
RUN go build -o main

ENTRYPOINT ["./scripts/startup.sh"]