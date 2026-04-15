#include <Arduino.h>
#include <ESP8266WiFi.h>

#include "Config.h"
#include "HALClient.h"
#include "HALDebug.h"
#include "RGBLed.h"
#include "Relay.h"
#include "Secrets.h"
#include "Switch.h"

unsigned long lastPoll = 0;

HALClient client(HAL_SERVER_URL, HAL_API_KEY);
bool relayState = false;

// setup() runs once on boot.
void setup() {
  Serial.begin(BAUD_RATE);

  delay(2000);

  debugPrintln("Connecting to WiFi...");

  WiFi.persistent(false);
  WiFi.disconnect(true);

  delay(100);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
  }

  debugPrintln("WiFi Connected!");
}

// loop() runs repeatedly after setup() completes.
void loop() {
  if (millis() - lastPoll >= POLL_INTERVAL_MS) {
    lastPoll = millis();

    debugPrintln("Polling...");

    char currentMessage[19];
    snprintf(currentMessage, sizeof(currentMessage), "Current state: %s",
             relayState ? "on" : "off");

    debugPrintln(currentMessage);

    bool shouldToggle = client.poll(relayState);

    if (shouldToggle) {
      relayState = !relayState;

      char toggleMessage[31];
      snprintf(toggleMessage, sizeof(toggleMessage),
               "Toggling state, new state: %s", relayState ? "on" : "off");

      debugPrintln(toggleMessage);
    }
  }
}