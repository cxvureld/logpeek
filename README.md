# LogPeek — Real-time Log Analyzer

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

Terminal-based real-time log analyzer for Nginx and JSON logs. Built with Go + Bubble Tea.

## ✨ Features

- Real-time metrics: Requests, status codes, bandwidth
- Top IPs & Endpoints: Automatic aggregation and sorting
- Multiple formats: Nginx (combined), JSON (CloudWatch style)
- Beautiful TUI: Built with Bubble Tea
- Streaming: Processes logs line-by-line, minimal memory
- Cross-platform: Linux, macOS, WSL

## 🚀 Quick Start

```git clone https://github.com/cxvureld/logpeek.git```
```cd logpeek```

Run with nginx logs:
```go run cmd/logpeek/main.go --file=test.log --format=nginx```

Run with JSON logs:
```go run cmd/logpeek/main.go --file=test_json.log --format=json```

Pipe from stdin:
```tail -f /var/log/nginx/access.log | go run cmd/logpeek/main.go --format=nginx```

## 📋 Supported Formats

Nginx Combined: ```--format=nginx``` (Standard nginx access logs)
JSON: ```--format=json``` (Structured JSON logs like CloudWatch, Loki)

## 🤔 Why LogPeek?

When production goes down, you don't have time to set up Grafana or write complex grep pipelines. You need answers right now:

- Which endpoint is failing? Top Endpoints with error rates
- Who's flooding the server? Top IP Addresses
- How many 5xx errors? Real-time status code breakdown
- All in one terminal command, no dependencies.

## 📜 License

MIT License
