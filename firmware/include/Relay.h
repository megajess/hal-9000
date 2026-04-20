#pragma once

#include <Arduino.h>

#include "Config.h"

// Phase 2.2 — Relay driver
// Controls the isolated relay module on PIN_RELAY (active HIGH).
// Toggling flips the current state.

class Relay {
 public:
  void begin() { pinMode(PIN_RELAY, OUTPUT); }

  void toggle() {
    _state = !_state;

    digitalWrite(PIN_RELAY, _state ? HIGH : LOW);
  }

  void setOn() {
    if (!_state) {
      toggle();
    }
  }

  void setOff() {
    if (_state) {
      toggle();
    }
  }

  bool isOn() const { return _state; }

 private:
  bool _state = false;
};
