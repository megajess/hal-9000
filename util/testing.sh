#!/bin/bash

set -e

SERVER_URL="http://localhost:8080"
DEVICE_NAME="Living Room"
SECRETS_FILE="../firmware/include/Secrets.h"

register() {
  echo "Registering device: $DEVICE_NAME"

  RESPONSE=$(curl -s -X POST "$SERVER_URL/devices" \
    -H "Content-Type: application/json" \
    -d "{\"name\": \"$DEVICE_NAME\"}")

  echo "Response: $RESPONSE"

  API_KEY=$(echo "$RESPONSE" | jq -r '.api_key')

  if [ -z "$API_KEY" ] || [ "$API_KEY" = "null" ]; then
    echo "Error: failed to extract API key from response"
    exit 1
  fi

  echo "API Key:   $API_KEY"

  awk -v key="$API_KEY" '
/^#define HAL_API_KEY/ {
  if (/\\$/) getline
  print "#define HAL_API_KEY    \"" key "\""
  next
}
{ print }
' "$SECRETS_FILE" > /tmp/hal_secrets_tmp && mv /tmp/hal_secrets_tmp "$SECRETS_FILE"

  echo "Secrets.h updated successfully"
}

set_state() {
  DESIRED_STATE=$1
  STATE_VALUE=0

  if [ "$DESIRED_STATE" = "on" ]; then
    STATE_VALUE=1
  fi

  API_KEY=$(grep '^#define HAL_API_KEY' "$SECRETS_FILE" | sed 's/.*"\(.*\)".*/\1/')

  if [ -z "$API_KEY" ]; then
    echo "Error: HAL_API_KEY not found in Secrets.h"
    exit 1
  fi

  echo "Setting device $API_KEY to $DESIRED_STATE"

  curl -X PATCH "$SERVER_URL/devices/state?state=$STATE_VALUE" \
  -H "X-API-Key: $API_KEY"

  echo ""
  echo "Done"
}

case "$1" in
  register) register ;;
  setOn)    set_state "on" ;;
  setOff)   set_state "off" ;;
  *)        echo "Usage: $0 {register|setOn|setOff}" ; exit 1 ;;
esac
