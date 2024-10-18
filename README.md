# Magnify 🔎
> Simple tool to summarize & show prometheus metrics via JSON api

![logo](logo.png)

## Usage
```bash
Usage of magnify:
  -c string
        yaml file config path (default "config.yaml")
```

## Example configuration
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