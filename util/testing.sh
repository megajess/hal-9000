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
  DEVICE_ID=$(echo "$RESPONSE" | jq -r '.id')

  if [ -z "$API_KEY" ] || [ "$API_KEY" = "null" ]; then
    echo "Error: failed to extract API key from response"
    exit 1
  fi

  echo "Device ID: $DEVICE_ID"
  echo "API Key:   $API_KEY"

  awk -v key="$API_KEY" '
/^#define HAL_API_KEY/ {
  if (/\\$/) getline
  print "#define HAL_API_KEY    \"" key "\""
  next
}
{ print }
' "$SECRETS_FILE" > /tmp/hal_secrets_tmp && mv /tmp/hal_secrets_tmp "$SECRETS_FILE"

  awk -v id="$DEVICE_ID" '
/^#define HAL_DEVICE_ID/ {
  if (/\\$/) getline
  print "#define HAL_DEVICE_ID  \"" id "\""
  next
}
{ print }
' "$SECRETS_FILE" > /tmp/hal_secrets_tmp && mv /tmp/hal_secrets_tmp "$SECRETS_FILE"

  echo "Secrets.h updated successfully"
}

set_state() {
  DESIRED_STATE=$1

  DEVICE_ID=$(grep '^#define HAL_DEVICE_ID' "$SECRETS_FILE" | sed 's/.*"\(.*\)".*/\1/')

  if [ -z "$DEVICE_ID" ]; then
    echo "Error: HAL_DEVICE_ID not found in Secrets.h"
    exit 1
  fi

  echo "Setting device $DEVICE_ID to $DESIRED_STATE"

  curl -s -X PUT "$SERVER_URL/devices/$DEVICE_ID" \
    -H "Content-Type: application/json" \
    -d "{\"desired_state\": \"$DESIRED_STATE\"}"

  echo ""
  echo "Done"
}

case "$1" in
  register) register ;;
  setOn)    set_state "on" ;;
  setOff)   set_state "off" ;;
  *)        echo "Usage: $0 {register|setOn|setOff}" ; exit 1 ;;
esac
