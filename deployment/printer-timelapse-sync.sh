#!/bin/bash
# Syncs timelapse videos from the Bambu printer's FTPS server to local storage.
# Runs hourly via systemd timer.
#
# Reads credentials from /opt/3d-printer/.env.secrets:
#   PRINTER_FTP_HOST, PRINTER_FTP_USER, PRINTER_FTP_PASSWORD

REMOTE_DIR='/timelapse'
LOCAL_DIR='/var/www/printer-timelapses'

if [ -z "$PRINTER_FTP_HOST" ] || [ -z "$PRINTER_FTP_USER" ] || [ -z "$PRINTER_FTP_PASSWORD" ]; then
  echo "ERROR: Missing PRINTER_FTP_HOST, PRINTER_FTP_USER, or PRINTER_FTP_PASSWORD"
  exit 1
fi

lftp -e "set ssl:verify-certificate no; set ftp:ssl-allow true; mirror -c --use-pget-n=10 --delete $REMOTE_DIR $LOCAL_DIR; quit" -p 990 -u "$PRINTER_FTP_USER","$PRINTER_FTP_PASSWORD" "$PRINTER_FTP_HOST"
