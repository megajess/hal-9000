#pragma once

#include <Arduino.h>

#include "Config.h"
#include "DeviceState.h"
#include "HALDebug.h"

// Phase 2.2 — RGB LED driver
// Controls the RGB LED to indicate device mode:
//   solid green     = normal
//   solid blue      = pairing
//   green blinking  = connecting to WiFi
//   red green blink = server unreachable
//   red blue blink  = WiFi offline

enum class PatternBlinkPhase { First, Second, Off };
enum class LEDColor { Red, Green, Blue };

class RGBLed {
 public:
  void begin() {
    pinMode(PIN_LED_R, OUTPUT);
    pinMode(PIN_LED_G, OUTPUT);
    pinMode(PIN_LED_B, OUTPUT);

    _lastBlinkTime = millis();
    _lastPatternBlinkTime = millis();

    off();
  }

  void indicateFor(DeviceState deviceState) {
    if (deviceState != _lastState) {
      _lastState = deviceState;
      _patternBlinkPhase = PatternBlinkPhase::First;
    }

    switch (deviceState) {
      case DeviceState::Connecting:
        indicateConnecting();
        break;
      case DeviceState::Connected:
        indicateOperational();
        break;
      case DeviceState::WifiOffline:
        indicateWifiOffline();
        break;
      case DeviceState::ServerUnreachable:
        indicateServerUnreachable();
        break;
    }
  }

  void off() {
    digitalWrite(PIN_LED_R, LOW);
    digitalWrite(PIN_LED_G, LOW);
    digitalWrite(PIN_LED_B, LOW);
  }

 private:
  // -- Properties --

  uint8_t _redState = LOW;
  uint8_t _greenState = LOW;
  uint8_t _blueState = LOW;
  bool _blinkOn = true;
  unsigned long _lastBlinkTime = 0;
  unsigned long _lastPatternBlinkTime = 0;
  unsigned long _patternBlinkInterval = 200;

  PatternBlinkPhase _patternBlinkPhase = PatternBlinkPhase::First;
  DeviceState _lastState = DeviceState::Connecting;

  // -- Methods --

  void setRed() {
    off();
    digitalWrite(PIN_LED_R, HIGH);
  }

  void setGreen() {
    off();
    digitalWrite(PIN_LED_G, HIGH);
  }

  void setBlue() {
    off();
    digitalWrite(PIN_LED_B, HIGH);
  }

  void indicateConnecting() { oneColorBlinkNonBlocking(LEDColor::Green); }

  void indicateOperational() { setGreen(); }

  void indicateWifiOffline() {
    patternBlinkNonBlocking(LEDColor::Red, LEDColor::Blue);
  }

  void indicateServerUnreachable() {
    patternBlinkNonBlocking(LEDColor::Red, LEDColor::Green);
  }

  void oneColorBlinkNonBlocking(LEDColor color) {
    if (millis() - _lastBlinkTime >= 250) {
      _lastBlinkTime = millis();
      _blinkOn = !_blinkOn;

      if (_blinkOn) {
        setLEDFor(color);
      } else {
        off();
      }
    }
  }

  void patternBlinkNonBlocking(LEDColor color1, LEDColor color2) {
    if (millis() - _lastPatternBlinkTime >= _patternBlinkInterval) {
      _lastPatternBlinkTime = millis();
      _patternBlinkInterval = 200;

      switch (_patternBlinkPhase) {
        case PatternBlinkPhase::First:
          setLEDFor(color1);
          _patternBlinkPhase = PatternBlinkPhase::Second;
          break;
        case PatternBlinkPhase::Second:
          setLEDFor(color2);
          _patternBlinkPhase = PatternBlinkPhase::Off;
          break;
        case PatternBlinkPhase::Off:
          off();
          _patternBlinkInterval = 500;
          _patternBlinkPhase = PatternBlinkPhase::First;
          break;
      }
    }
  }

  void setLEDFor(LEDColor color) {
    switch (color) {
      case LEDColor::Red:
        setRed();
        break;
      case LEDColor::Green:
        setGreen();
        break;
      case LEDColor::Blue:
        setBlue();
        break;
    }
  }
};
