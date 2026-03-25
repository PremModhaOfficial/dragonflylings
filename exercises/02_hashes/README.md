# Module 02: Hashes — "The Row"

## Mental Model

If a String is a single cell in a spreadsheet, a Hash is an entire row. One key, many fields. You wouldn't create 10 separate keys for a user's name, email, age — you'd use one hash. It's the difference between 10 phone calls and one phone call with 10 questions.

```
String approach (scattered):          Hash approach (one key):
─────────────────────────────         ─────────────────────────────────
"user:42:name"   = "alice"            "user:42"
"user:42:email"  = "a@ex.com"            name    = "alice"
"user:42:age"    = "29"                  email   = "a@ex.com"
"user:42:role"   = "admin"               age     = "29"
                                         role    = "admin"

To delete user:                       To delete user:
DEL user:42:name                      DEL user:42
DEL user:42:email                     (one command)
DEL user:42:age
DEL user:42:role
```

## Predict Before Starting

Before writing any code, answer in your head:
1. When would you choose multiple String keys over a Hash for the same data?
2. What does HGETALL return for a hash with 100 fields?
3. Can a hash field hold a number? Can you HINCRBY it?
4. Does deleting a hash field (HDEL) affect other fields?

Write your predictions in `feynman/gap_notebook.md`.

## Key Concepts

| Command | Description |
|---------|-------------|
| HSET key field value [field value ...] | Set one or more hash fields |
| HGET key field | Get a single field |
| HMGET key field [field ...] | Get multiple fields (nil for missing) |
| HGETALL key | Get all fields and values |
| HINCRBY key field n | Atomically increment a numeric field |
| HDEL key field | Delete a field |
| HLEN key | Number of fields in hash |

## Exercises

- **hashes1**: HSET/HGET — model a user profile as a hash (not scattered string keys)
- **hashes2**: HMGET — batch field reads in one round trip vs. N round trips
- **hashes3**: HINCRBY — counters inside hashes (page view tracker per page)
- **hashes4**: Hash vs. multiple Strings — understand when to use each

## Before You Start

```bash
# Experiment with hashes
redis-cli -p 6380 HSET user:1 name alice email alice@example.com age 29
redis-cli -p 6380 HGET user:1 name
redis-cli -p 6380 HGETALL user:1
redis-cli -p 6380 TYPE user:1   # "hash"
```
