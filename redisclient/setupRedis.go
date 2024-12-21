package redisclient

import (
  "github.com/redis/go-redis/v9"
)

func SetupRedisCaching() *redis.Client {
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",            // No Password
    DB: 0,                   // Use Default Db
    Protocol: 2,             // connection protocol
  })

  return client
}
