#pragma once

#include <Arduino.h>

inline void debugPrintln(const char* message) {
#ifdef DEBUG
  Serial.println(message);
#endif
}
