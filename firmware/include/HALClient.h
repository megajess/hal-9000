#pragma once

#include <Arduino.h>
#include <ESP8266HTTPClient.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h>

#include "HALDebug.h"
#include "Secrets.h"

// HALClient handles all communication with the HAL server.
// Polls the server with current device state and returns
// whether the relay should toggle.

enum class PollResult { Toggle, NoChange, Failure };

class HALClient {
 public:
  HALClient(const char* serverUrl, const char* apiKey)
      : _serverUrl(serverUrl), _apiKey(apiKey) {}

  PollResult poll(bool currentState) {
    HTTPClient http;

    char url[128];
    snprintf(url, sizeof(url), "%s/poll?state=%d", _serverUrl,
             currentState ? 1 : 0);

    http.begin(_wifiClient, url);
    http.addHeader("X-API-Key", _apiKey);

    int statusCode = http.GET();

    if (statusCode == 200) {
      String response = http.getString();

      char debug[64];
      snprintf(debug, sizeof(debug), "Poll status: %d response: %s", statusCode,
               response.c_str());

      debugPrintln(debug);

      http.end();

      return response == "1" ? PollResult::Toggle : PollResult::NoChange;
    } else {
      char debug[64];
      snprintf(debug, sizeof(debug), "Poll failed, status: %d", statusCode);
      debugPrintln(debug);

      http.end();

      return PollResult::Failure;
    }
  }

  bool setDesiredState(bool state) {
    HTTPClient http;

    char url[64];
    snprintf(url, sizeof(url), "%s/devices/state?state=%d", _serverUrl,
             state ? 1 : 0);

    http.begin(_wifiClient, url);
    http.addHeader("X-API-Key", _apiKey);

    int statusCode = http.sendRequest("PATCH", "");

    if (statusCode == 200) {
      http.end();

      return true;
    } else {
      char debug[64];
      snprintf(debug, sizeof(debug), "Setting desired state failed, status: %d",
               statusCode);
      debugPrintln(debug);

      http.end();

      return false;
    }
  }

 private:
  const char* _serverUrl;
  const char* _apiKey;
  WiFiClient _wifiClient;
};
