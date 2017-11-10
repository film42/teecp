Teecp
=====

Teecp is a transparent tcp proxy that allows you to "tee" (or replicate) your traffic to additional servers.

### Installing

```
$ go get github.com/film42/teecp
```

### Usage

```
Usage of ./teecp:
  -config string
    	Path to the teecp config. (default "config.json")
  -debug
    	Enable debug logging.
```

Here's an example config file:

```json
{
  "bind": "localhost:10000",
  "proxy": "localhost:6379",
  "tees": [
    "corvus_session:6380"
  ]
}
```

This tool can be useful when migrating from redis sentinel to redis cluster without losing data. We use redis to store
session data, and needed to find a way to migrate to redis cluster without resetting user sessions. We used teecp and
[corvus](https://github.com/eleme/corvus) to replicate our redis traffic to both standalone redis and redis cluster.
Once both redis and redis cluster were consistent, we were able to use teecp to proxy client reads and writes to redis
cluster, while still replicating to standalone redis to ensure consistency until our entire infrastructure was on redis
cluster. After that, we removed teecp and connected directly to redis cluster.

### License

MIT License
