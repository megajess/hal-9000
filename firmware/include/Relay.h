#pragma once

// Phase 2.2 — Relay driver
// Controls the isolated relay module on PIN_RELAY (active HIGH).
// Toggling flips the current state.

class Relay {
public:
    void begin();
    void toggle();
    bool isOn() const;

private:
    bool _state = false;
};
