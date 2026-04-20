#pragma once

#include <Arduino.h>

#include "Config.h"

// Phase 2.2 — Physical switch input
// Reads the momentary SPDT center off switch.
// UP   (PIN_SWITCH_ON)  → always on
// DOWN (PIN_SWITCH_OFF) → always off
// CENTER                → no change

enum class SwitchPosition { Up, Down, Center };

class Switch {
 public:
  void begin() {
    pinMode(PIN_SWITCH_ON, INPUT_PULLUP);
    pinMode(PIN_SWITCH_OFF, INPUT_PULLUP);
  }

  SwitchPosition readDebounced() {
    SwitchPosition reading = read();

    if (reading != _lastReading) {
      _lastDebounceTime = millis();
    }

    if (millis() - _lastDebounceTime > DEBOUNCE_DELAY_MS) {
      _stablePosition = reading;
    }

    _lastReading = reading;
    return _stablePosition;
  }

 private:
  SwitchPosition _lastReading = SwitchPosition::Center;
  SwitchPosition _stablePosition = SwitchPosition::Center;
  unsigned long _lastDebounceTime = 0;

  SwitchPosition read() const {
    if (digitalRead(PIN_SWITCH_ON) == LOW) {
      return SwitchPosition::Up;
    } else if (digitalRead(PIN_SWITCH_OFF) == LOW) {
      return SwitchPosition::Down;
    } else {
      return SwitchPosition::Center;
    }
  }
};
