#pragma once

// --- Pin Definitions ---

// Relay — active HIGH, switches 120V load
#define PIN_RELAY 5  // D1

// Physical switch — momentary SPDT center off, low voltage
#define PIN_SWITCH_OFF 4  // D2 — switch UP, always off
#define PIN_SWITCH_ON 0   // D3 — switch DOWN, always on

// RGB LED — common cathode, driven via PWM
#define PIN_LED_R 14  // D5
#define PIN_LED_G 12  // D6
#define PIN_LED_B 13  // D7

// Pairing button — must have external pull down, affects boot mode
#define PIN_PAIRING 15  // D8

// --- Timing ---
#define POLL_INTERVAL_MS 2000  // how often to poll the server (ms)
#define WIFI_TIMEOUT_MS 10000  // how long to wait for WiFi on boot (ms)

// --- Serial ---
#define BAUD_RATE 115200

// --- Network ---
#define STATIC_IP "192.168.50.78"
#define GATEWAY_IP "192.168.50.1"
#define SUBNET_MASK "255.255.255.0"
#define DNS_IP "192.168.50.1"

// --- Hardware ---
#define DEBOUNCE_DELAY_MS 50

// --- Connectivity ---
#define MAX_POLL_FAILURES 3
constexpr unsigned long RECONNECT_BACKOFF_MS[] = {5000, 10000, 20000, 30000};
constexpr int RECONNECT_BACKOFF_LEN =
    sizeof(RECONNECT_BACKOFF_MS) / sizeof(RECONNECT_BACKOFF_MS[0]);