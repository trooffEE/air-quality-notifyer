FROM golang:1.24.3

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN chmod +x ./scripts/startup.sh
RUN chmod +x ./scripts/shutdown.sh

ENTRYPOINT ["./scripts/startup.sh"]