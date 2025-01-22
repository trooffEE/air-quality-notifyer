#!/bin/sh

cp ./scripts/pre-push.sh .git/hooks/pre-push
chmod 755 .git/hooks/pre-push