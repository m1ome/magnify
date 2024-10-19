# Magnify 🔎
> Simple tool to summarize & show prometheus metrics via JSON api

![logo](logo.png)

## About
**Magnify** is a simple tool to provide a JSON API based metrics resolution on top of [Prometheus](https://prometheus.io/) it heavily relies on [Expr](https://expr-lang.org/docs/getting-started) for metrics transformation.

## Usage
```bash
Usage of magnify:
  -c string
        yaml file config path (default "config.yaml")
```

## Example configuration
Please look at [Expr](https://expr-lang.org/docs/getting-started) documentation for expressions.

```yaml
prometheus:
  addr: "http://localhost:4000"
expressions:
  - name: goroutines
    query: sum(rate(go_gc_duration_seconds{app_kubernetes_io_name="game-server"}[5m])) by (service)
    expr: "float(query_result[0].Value.String()) >= 0.0 ? 'operational': 'error'"
http:
  addr: ":9999"
```

## Endpoints

### Cache
Method: **GET**
Path: **/**

Response:
```json
{"goroutines": "operational"}
```

### Key
Method: **GET**
Path: **/{name}**

Response:
```json
"operational"
```