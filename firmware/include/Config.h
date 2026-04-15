#pragma once

// --- Pin Definitions ---

// Relay — active HIGH, switches 120V load
#define PIN_RELAY       5   // D1

// Physical switch — momentary SPDT center off, low voltage
#define PIN_SWITCH_ON   4   // D2 — switch UP, always on
#define PIN_SWITCH_OFF  16  // D0 — switch DOWN, always off

// RGB LED — common cathode, driven via PWM
#define PIN_LED_R       14  // D5
#define PIN_LED_G       12  // D6
#define PIN_LED_B       13  // D7

// Pairing button — must have external pull down, affects boot mode
#define PIN_PAIRING     15  // D8

// --- Timing ---
#define POLL_INTERVAL_MS  2000  // how often to poll the server (ms)
#define WIFI_TIMEOUT_MS   10000 // how long to wait for WiFi on boot (ms)

// --- Serial ---
#define BAUD_RATE  115200
