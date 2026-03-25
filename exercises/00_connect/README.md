# Module 00: Connect — "The Handshake"

## Mental Model

A Redis connection is like a phone call. You dial (TCP), someone answers (RESP handshake), you say "hello?" (PING), they say "I'm here" (PONG). Only then do you talk.

Just creating a client object does **not** establish a connection. go-redis is lazy: the TCP connection only opens on the first command. This means your app can "appear" connected while Dragonfly is completely unreachable.

```
Your App          Network         Dragonfly
   |                |                 |
   |--NewClient()   |                 |  <- no network call yet
   |                |                 |
   |--Ping()------->|---------------->|  <- TCP SYN
   |                |                 |--TCP SYN-ACK-->
   |                |<-"PONG"---------|
   |<-"PONG"--------|                 |
```

## Predict Before Starting

Before writing any code, answer in your head:
1. What does `redis.NewClient()` actually do — does it open a connection?
2. What does PING return? Why would that command exist in a database protocol?
3. Dragonfly runs on port 6380 in this project. Why not 6379 (the Redis default)?
4. If you configure a pool of 10 connections, when are those connections established?

Write your predictions in `feynman/gap_notebook.md` before looking at the exercises.

## Key Concepts

| Concept | Description |
|---------|-------------|
| Lazy connection | `redis.NewClient()` creates a client struct, not a TCP connection |
| PING/PONG | Health check command — verifies the full protocol handshake |
| DialTimeout | How long to wait for TCP connection before giving up |
| PoolSize | Max concurrent connections to Dragonfly |
| MinIdleConns | Pre-warmed idle connections ready for immediate use |

## Exercises

- **connect1**: Fix the wrong port and wrong method — your first PING
- **connect2**: Handle unreachable Dragonfly — what should your code do?
- **connect3**: Configure a connection pool and observe its behavior

## Before You Start

```bash
# Make sure Dragonfly is running
docker compose up -d

# Verify it's alive
redis-cli -p 6380 PING
# Should return: PONG
```
