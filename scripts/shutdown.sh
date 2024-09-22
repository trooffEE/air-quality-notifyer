#!/bin/bash

if [[ ! -z $DEVELOPMENT ]]; then
  echo "DEVELOPMENT: Setting Webhook back"
  curl "https://api.telegram.org/bot${TELEGRAM_SECRET}/setWebhook?url=https://${WEBHOOK_HOST}/webhook${TELEGRAM_SECRET}"
fi
