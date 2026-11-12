# ⚡ NanoKV: Distributed Key-Value Store

> A high-performance, fault-tolerant distributed key-value database built from scratch in Go. Features **Sharding**, **Write-Ahead Logging (WAL)**, and **Consistent Hashing**.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Distributed-orange?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production_Ready-success?style=for-the-badge)

## 📸 Live System Demo

![NanoKV Architecture Demo](NanoKV.jpeg)
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
```

---

## ⚡ Custom Development, Commercial Support & Deployment

Need a custom **distributed backend service**, **low-latency key-value / caching layer**, **scalable REST API (Node.js / TypeScript / Java Spring Boot / Go)**, or **multi-container Docker & Linux server deployment** built for your infrastructure?

I offer fast-turnaround software engineering services and custom proposals on Fiverr:

[![Hire on Fiverr - Scalable REST APIs & Distributed Backends](https://img.shields.io/badge/Fiverr-Scalable_REST_APIs_%26_Backend_Architecture_(From_$35)-00b22d?style=for-the-badge&logo=fiverr&logoColor=white)](https://www.fiverr.com/thepimen/develop-scalable-rest-apis-with-nodejs-typescript-or-spring-boot)
[![Hire on Fiverr - Docker, Microservices & Linux Servers](https://img.shields.io/badge/Fiverr-Docker,_Microservices_%26_Linux_Servers_(From_$35)-1dbf73?style=for-the-badge&logo=fiverr&logoColor=white)](https://www.fiverr.com/thepimen/dockerize-applications-and-build-microservices-on-linux-servers)

* **24-Hour Express Delivery available** for urgent backend APIs, caching bottlenecks, and Docker/Linux deployments.
* **Direct Engineering Portfolio & Services:** [luislazaro-dev.netlify.app/solutions](https://luislazaro-dev.netlify.app/solutions)

