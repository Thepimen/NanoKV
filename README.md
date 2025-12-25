# ⚡ NanoKV: Distributed Key-Value Store

> A high-performance, fault-tolerant distributed key-value database built from scratch in Go. Features **Sharding**, **Write-Ahead Logging (WAL)**, and **Consistent Hashing**.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Distributed-orange?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production_Ready-success?style=for-the-badge)

## 📸 Live System Demo

![NanoKV Architecture Demo](demo-architecture.png)
*Real-time cluster operation showing the Proxy (right) routing requests via Consistent Hashing to specific Shards (left), with persistence logs active.*

---

## 🧠 System Architecture

NanoKV is designed as a distributed system where data is partitioned across multiple nodes to ensure horizontal scalability.

```mermaid
graph TD
    Client[Client / CLI] -->|HTTP POST /set| Proxy[LB / Proxy Node :9000]
    Proxy -->|Consistent Hashing| NodeA[Shard 0 :8080]
    Proxy -->|Consistent Hashing| NodeB[Shard 1 :8081]
    Proxy -->|Consistent Hashing| NodeC[Shard 2 :8082]
    NodeA --> Disk[(WAL Log)]
    NodeB --> Disk
    NodeC --> Disk
