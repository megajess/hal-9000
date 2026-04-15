#pragma once

// Phase 2.2 — RGB LED driver
// Controls the RGB LED to indicate device mode:
//   green  = normal
//   blue   = pairing
//   red    = debug

class RGBLed {
public:
    void begin();
    void setGreen();
    void setBlue();
    void setRed();
    void off();
};
