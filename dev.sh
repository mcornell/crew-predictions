#!/usr/bin/env bash
set -e

# Kill anything already on these ports (emulators + Go server)
kill $(lsof -ti :8080 :8081 :9099 :4000 :4400 :4500) 2>/dev/null || true

# Start Firestore + Auth emulators in background
firebase emulators:start --only firestore,auth &
EMULATOR_PID=$!

# Wait for emulators to be ready
sleep 5

# Start Go server
FIRESTORE_EMULATOR_HOST=localhost:8081 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
GOOGLE_CLOUD_PROJECT=crew-predictions \
/usr/local/go/bin/go run ./cmd/server

# Clean up emulators on exit
kill $EMULATOR_PID 2>/dev/null || true
