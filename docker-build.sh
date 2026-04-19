#!/usr/bin/env bash
# Build and optionally smoke-test the production Docker image locally.
# Only needed when the Dockerfile changes — for code-only changes, push
# directly and let Cloud Build handle it.
#
# Usage:
#   ./docker-build.sh          # build only
#   ./docker-build.sh --run    # build + start container on :8080

set -e

IMAGE="crew-predictions:local"

# Ensure Docker daemon is running
if ! docker info >/dev/null 2>&1; then
  echo "Docker daemon not running. Starting Docker Desktop..."
  systemctl --user start docker-desktop
  echo "Waiting for Docker..."
  until docker info >/dev/null 2>&1; do sleep 2; done
  echo "Docker ready."
fi

echo "Building $IMAGE..."
docker build -t "$IMAGE" .
echo "Build complete."

if [[ "$1" == "--run" ]]; then
  # Load env vars from .env (skip comments and blank lines)
  ENV_ARGS=()
  if [[ -f .env ]]; then
    while IFS= read -r line; do
      [[ "$line" =~ ^#|^$ ]] && continue
      ENV_ARGS+=(-e "$line")
    done < .env
  fi

  echo "Starting container on http://localhost:8080 (Ctrl-C to stop)..."
  docker run --rm -p 8080:8080 "${ENV_ARGS[@]}" "$IMAGE"
fi
