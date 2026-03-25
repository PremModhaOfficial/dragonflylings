## Hint 1

`JSON.SET` expects the value to be a **JSON string**, not a Go struct. You need to marshal the Profile first:

```go
data, err := json.Marshal(p)
// data is now []byte: {"name":"Alice","age":30,...}
// pass string(data) or data to JSON.SET
```

## Hint 2

`JSON.GET key $` returns a **JSON array** wrapping the result: `[{"name":"Alice",...}]`

To unmarshal a Profile, you need to unwrap the array:
```go
var profiles []Profile
json.Unmarshal([]byte(result), &profiles)
return profiles[0], nil
```

## Hint 3

Complete fix for both functions:

```go
func SetProfile(ctx context.Context, client *redis.Client, key string, p Profile) error {
    data, err := json.Marshal(p)
    if err != nil {
        return err
    }
    return client.Do(ctx, "JSON.SET", key, "$", string(data)).Err()
}

func GetProfile(ctx context.Context, client *redis.Client, key string) (Profile, error) {
    result, err := client.Do(ctx, "JSON.GET", key, "$").Text()
    if err != nil {
        return Profile{}, err
    }
    var profiles []Profile
    if err := json.Unmarshal([]byte(result), &profiles); err != nil {
        return Profile{}, err
    }
    if len(profiles) == 0 {
        return Profile{}, redis.Nil
    }
    return profiles[0], nil
}
```
