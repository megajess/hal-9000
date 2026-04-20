#include <Arduino.h>
#include <ESP8266WiFi.h>

#include "Config.h"
#include "DeviceState.h"
#include "HALClient.h"
#include "HALDebug.h"
#include "RGBLed.h"
#include "Relay.h"
#include "Secrets.h"
#include "Switch.h"

unsigned long lastPoll = 0;
int consecutivePollFailures = 0;
bool switchUsedOffline = false;
unsigned long lastReconnectAttempt = 0;
int reconnectBackoffIndex = 0;

DeviceState deviceState = DeviceState::Connecting;
HALClient client(HAL_SERVER_URL, HAL_API_KEY);
Relay relay = Relay();
Switch physicalSwitch = Switch();
SwitchPosition lastPosition = SwitchPosition::Center;
RGBLed led = RGBLed();

// -- Function prototypes --
void handleSwitch();
void handleLED();
void handleConnecting();
void handleConnected();
void handleWifiOffline();
void handleServerUnreachable();
void attemptReconnect();
bool handlePoll();
void debug_printToggleMessage();
void debug_printCurrentStateMessage();

// setup() runs once on boot.
void setup() {
  Serial.begin(BAUD_RATE);

  physicalSwitch.begin();
  relay.begin();
  led.begin();

  delay(2000);

  WiFi.mode(WIFI_STA);
  WiFi.persistent(false);
  WiFi.disconnect(true);

  delay(500);

  IPAddress ip, gateway, subnet, dns;
  ip.fromString(STATIC_IP);
  gateway.fromString(GATEWAY_IP);
  subnet.fromString(SUBNET_MASK);
  dns.fromString(DNS_IP);

  WiFi.config(ip, gateway, subnet, dns);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  debugPrintln("Connecting to WiFi...");
}

// loop() runs repeatedly after setup() completes.
void loop() {
  handleSwitch();
  handleLED();

  switch (deviceState) {
    case DeviceState::Connecting:
      handleConnecting();
      break;
    case DeviceState::Connected:
      handleConnected();
      break;
    case DeviceState::WifiOffline:
      handleWifiOffline();
      break;
    case DeviceState::ServerUnreachable:
      handleServerUnreachable();
      break;
  }
}

void handleSwitch() {
  SwitchPosition position = physicalSwitch.readDebounced();

  if (position != lastPosition) {
    switchUsedOffline = deviceState != DeviceState::Connected;

    if (position == SwitchPosition::Up) {
      if (!relay.isOn()) relay.toggle();

      if (deviceState == DeviceState::Connected) {
        client.setDesiredState(true);
      }
    } else if (position == SwitchPosition::Down) {
      if (relay.isOn()) relay.toggle();

      if (deviceState == DeviceState::Connected) {
        client.setDesiredState(false);
      }
    }

    lastPosition = position;
  }
}

void handleLED() { led.indicateFor(deviceState); }

void handleConnecting() { attemptReconnect(); }

void handleConnected() {
  if (switchUsedOffline) {
    bool desiredStateUpdated = client.setDesiredState(relay.isOn());

    if (desiredStateUpdated) {
      switchUsedOffline = false;
    }
  }

  if (WiFi.status() != WL_CONNECTED) {
    deviceState = DeviceState::WifiOffline;
    consecutivePollFailures = 0;

    return;
  }

  handlePoll();
}

void handleWifiOffline() { attemptReconnect(); }

void handleServerUnreachable() {
  if (WiFi.status() != WL_CONNECTED) {
    deviceState = DeviceState::WifiOffline;
    return;
  }
  if (handlePoll()) {
    deviceState = DeviceState::Connected;
  }
}

void attemptReconnect() {
  if (millis() - lastReconnectAttempt >=
      RECONNECT_BACKOFF_MS[reconnectBackoffIndex]) {
    lastReconnectAttempt = millis();

    if (reconnectBackoffIndex < RECONNECT_BACKOFF_LEN - 1) {
      reconnectBackoffIndex++;
    }

    debugPrintln("Attempting reconnect...");

    WiFi.disconnect(true);
    delay(100);
    WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  }

  if (WiFi.status() == WL_CONNECTED) {
    reconnectBackoffIndex = 0;
    deviceState = DeviceState::Connected;
  }
}

// Returns true to indicate if the device state should change back to connected,
// if it is not already in connected state.
bool handlePoll() {
  if (millis() - lastPoll >= POLL_INTERVAL_MS) {
    lastPoll = millis();

    debugPrintln("Polling...");
    debug_printCurrentStateMessage();

    PollResult result = client.poll(relay.isOn());

    switch (result) {
      case PollResult::Toggle:
        relay.toggle();

        debug_printToggleMessage();
      // Fallthrough to reset consecutivePollFailures, and return notice of
      // successful poll
      [[fallthrough]]
      case PollResult::NoChange:
        consecutivePollFailures = 0;

        return true;
      case PollResult::Failure:
        if (++consecutivePollFailures >= MAX_POLL_FAILURES) {
          deviceState = DeviceState::ServerUnreachable;
        }

        return false;
    }
  }

  return false;
}

void debug_printToggleMessage() {
#ifdef DEBUG
  char toggleMessage[31];
  snprintf(toggleMessage, sizeof(toggleMessage),
           "Toggling state, new state: %s", relay.isOn() ? "on" : "off");

  debugPrintln(toggleMessage);
#endif
}

void debug_printCurrentStateMessage() {
#ifdef DEBUG
  char currentMessage[19];
  snprintf(currentMessage, sizeof(currentMessage), "Current state: %s",
           relay.isOn() ? "on" : "off");

  debugPrintln(currentMessage);
#endif
}