#pragma once

// Phase 2.2 — Physical switch input
// Reads the momentary SPDT center off switch.
// UP   (PIN_SWITCH_ON)  → always on
// DOWN (PIN_SWITCH_OFF) → always off
// CENTER                → no change

enum class SwitchPosition {
    Up,
    Down,
    Center
};

class Switch {
public:
    void begin();
    SwitchPosition read() const;
};
