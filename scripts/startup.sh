#!/bin/bash

if [[ ! -z $DEVELOPMENT ]]; then
  echo "DEVELOPMENT: Deleting Webhook"
  curl "https://api.telegram.org/bot${TELEGRAM_SECRET}/deleteWebhook?url=https://${WEBHOOK_HOST}/webhook${TELEGRAM_SECRET}"

  make run

  exit 0
fi

make build

./main &
PID=$!
echo "PID: $PID"

handle_sigterm() {
    /app/scripts/shutdown.sh
    exit 0
}

trap 'handle_sigterm' SIGTERM SIGINT
wait $PID