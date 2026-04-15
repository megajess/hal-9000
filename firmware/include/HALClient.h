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

class HALClient {
 public:
  HALClient(const char* serverUrl, const char* apiKey)
      : _serverUrl(serverUrl), _apiKey(apiKey) {}

  bool poll(bool currentState) {
    HTTPClient http;

    char url[128];
    snprintf(url, sizeof(url), "%s/poll?state=%d", _serverUrl,
             currentState ? 1 : 0);

    http.begin(_wifiClient, url);
    http.addHeader("X-API-Key", _apiKey);

    int statusCode = http.GET();

    if (statusCode > 0) {
      String response = http.getString();

      char debug[64];
      snprintf(debug, sizeof(debug), "Poll status: %d response: %s", statusCode,
               response.c_str());

      debugPrintln(debug);

      http.end();

      return response == "1";
    } else {
      char debug[64];
      snprintf(debug, sizeof(debug), "Poll failed, status: %d", statusCode);
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
