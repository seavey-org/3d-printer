#!/bin/bash
set -euo pipefail

# Watchdog for ffmpeg printer stream
# Checks if stream.m3u8 is stale (>60s old) and restarts ffmpeg if so.
# Runs via systemd timer every 30 seconds.

STREAM_FILE="/var/www/printer-camera/live/stream.m3u8"
STALE_THRESHOLD=60
SERVICE="ffmpeg-printer-stream.service"

if [ ! -f "$STREAM_FILE" ]; then
    echo "$(date): stream.m3u8 missing, restarting $SERVICE"
    systemctl restart "$SERVICE"
    exit 0
fi

# Get file age in seconds
FILE_MTIME=$(stat -c %Y "$STREAM_FILE")
NOW=$(date +%s)
AGE=$((NOW - FILE_MTIME))

if [ "$AGE" -gt "$STALE_THRESHOLD" ]; then
    echo "$(date): stream.m3u8 is ${AGE}s old (threshold: ${STALE_THRESHOLD}s), restarting $SERVICE"
    systemctl restart "$SERVICE"
else
    echo "$(date): stream.m3u8 is ${AGE}s old, OK"
fi
