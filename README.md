# healthgate

> High-speed synthetic network, HTTP status, TCP reachability, and DNS resolution probe with **Prometheus exporter** written in **Go**.

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Prometheus](https://img.shields.io/badge/Metrics-Prometheus-E6522C?style=flat-square&logo=prometheus)](https://prometheus.io)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://docker.com)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

`#network-monitoring` `#synthetic-probe` `#prometheus-exporter` `#dns-probe` `#tcp-probe` `#http-probe` `#devops` `#golang`

---

## Features

- **HTTP/S Probing:** Measures status codes and roundtrip latency without following arbitrary redirects.
- **TCP Port Testing:** Verifies raw TCP socket connectivity (databases, caches, custom daemons).
- **DNS Resolution Latency:** Measures lookup durations across specific custom nameservers (e.g. 8.8.8.8, 1.1.1.1).
- **Prometheus Metrics:** Native histogram and gauge exports on `/metrics`.

## Quick Start

```bash
# Copy sample configuration
cp config.example.yaml config.yaml

# Run probe
go run main.go -config=config.yaml -listen=:9115
```

## Docker

```bash
docker build -t healthgate .
docker run -d -p 9115:9115 -v $(pwd)/config.yaml:/etc/healthgate/config.yaml healthgate
```
