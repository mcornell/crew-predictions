FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/* \
 && useradd -r -u 1001 -s /sbin/nologin app
WORKDIR /app
COPY server .
RUN chown -R app:app /app
USER app
EXPOSE 8080
CMD ["./server"]
