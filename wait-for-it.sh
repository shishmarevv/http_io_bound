#!/usr/bin/env bash
set -e

if [ $# -lt 2 ]; then
  echo "Usage: $0 host:port cmd [args...]"
  exit 1
fi

HOSTPORT=$1
shift

# Парсим хост и порт
IFS=":" read -r HOST PORT <<< "$HOSTPORT"

for i in {1..30}; do
  nc -z "$HOST" "$PORT" && break
  echo "[$(date +%H:%M:%S)] Waiting for $HOSTPORT... ($i/30)"
  sleep 1
done

exec "$@"