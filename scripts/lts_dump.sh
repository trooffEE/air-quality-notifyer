#!/usr/bin/env sh

find_lts_dump() {
  lts_dump=$(ls -t backup_*.sql 2>/dev/null | head -n 1)
  if [ -z "$lts_dump" ]; then
    echo "No dump file found."
    return 1
  fi
  return 0
}

applyDump() {
  if find_lts_dump; then
    echo "Applying dump: $lts_dump"
    docker cp "./data/dump/$lts_dump" airquality-db-container:/data/dump
    docker exec -it airquality-db-container sh -c "psql -U ${DB_USER} ${DB_NAME} < /data/dump/$lts_dump"
    docker exec -it airquality-db-container sh -c "rm /data/dump/$lts_dump"
  else
    echo "Error: No dump file found to apply."
    exit 1
  fi
}

applyDump
