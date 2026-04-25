FROM node:24-slim AS frontend
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY src/ ./src/
COPY index.html vite.config.ts tsconfig.json ./
RUN ./node_modules/.bin/vite build

FROM golang:1.26 AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN go build -o server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* \
 && useradd -r -u 1001 -s /sbin/nologin app
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=frontend /app/dist ./dist
RUN chown -R app:app /app
USER app
EXPOSE 8080
CMD ["./server"]
