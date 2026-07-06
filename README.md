# Radiohead Ticketing System

A high-performance, enterprise-grade ticketing backend designed to handle extreme concurrent traffic (Flash Sales). Built with Go, Redis, PostgreSQL, and Kafka.

## Architecture & Features

This system is engineered for zero-downtime, high-concurrency ticket booking, eliminating the need for separate scaling systems during traffic spikes. It relies on three core pillars:

1. **Virtual Waiting Room (Fail-fast via Redis)**: Sudden traffic spikes are absorbed by a Redis Sorted Set. Users are queued chronologically and admitted in batches, preventing the main application and database from being overwhelmed.
2. **Optimistic Locking (PostgreSQL)**: Once admitted, the ticket purchasing logic utilizes Optimistic Locking (`version` checks) directly at the database layer. This ensures absolute consistency and prevents double-booking (overbooking) race conditions without the heavy performance penalty of table/row locks.
3. **Asynchronous Event-Driven Processing (Kafka)**: Upon a successful purchase, the main transaction finishes instantly by publishing a `TicketSoldEvent` to Kafka. Downstream systems (e.g., Email, SMS, Analytics) consume these events asynchronously, ensuring the HTTP request resolves in milliseconds.

## Technology Stack
- **Backend:** Go (`net/http`)
- **Database:** PostgreSQL (with Optimistic Locking)
- **In-Memory Cache & Queue:** Redis
- **Message Broker:** Apache Kafka (KRaft mode)
- **Load Testing:** k6
- **Monitoring:** Prometheus & Grafana
- **Infrastructure:** Docker & Docker Compose

## Performance & Load Test Results

We battle-tested the system using [k6](https://k6.io/) simulating a massive "Flash Sale" event on a single developer machine (Docker on Windows, 8 Cores / 16 Threads, 32GB RAM).

### 🚀 Run 1: 5,000 Concurrent Virtual Users
*Simulating 5,000 users constantly hitting the system for 50 seconds.*

- **Total Requests:** `152,018`
- **Success Rate:** `100.00%` (0 dropped connections)
- **Throughput:** `2,984 RPS`
- **Median Latency:** `95ms` (50% of requests handled in under 1/10th of a second!)

### 💥 Run 2: "Breaking Point" at 50,000 Concurrent Virtual Users
*Simulating a massive surge to find the system's physical limits.*

- **Total Requests:** `173,224`
- **Throughput:** Maxed out Docker's networking allocation (`~3,220 RPS`)
- **Result:** `83%` of requests timed out (`dial: i/o timeout`).
- **Insight:** At 50,000 concurrent users, the host OS and Docker's Network Stack ran out of available sockets (TCP Backlog filled up). Crucially, the **Go Server did not crash**. It gracefully continued to serve up to the OS's networking physical limits, proving the architecture's incredible resilience.

## Getting Started

1. Clone the repository.
2. Ensure you have Docker installed.
3. Spin up the entire infrastructure (Go API, Postgres, Redis, Kafka, Prometheus, Grafana):
   ```bash
   docker-compose up -d --build
   ```
4. Access the API at `http://localhost:8080`.
