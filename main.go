package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/expr-lang/expr"
	"github.com/m1ome/magnify/pkg"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "yaml file config path")
	flag.Parse()
}

func main() {
	cfg, err := pkg.NewConfig(configPath)
	if err != nil {
		log.Fatalf("error initialzing config: %v", err)
	}

	c, err := api.NewClient(api.Config{
		Address: cfg.Prometheus.Addr,
	})
	if err != nil {
		log.Fatalf("error connecting to prometheus: %v", err)
	}
	v1api := v1.NewAPI(c)

	log.Print("scraping metrics on a startup")
	cache := pkg.NewCache(time.Duration(cfg.Cache.Expiration) * time.Second)

	scrapeMetrics(cfg, cache, v1api)

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			scrapeMetrics(cfg, cache, v1api)
			cache.Cleanup()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		m := cache.Copy()

		buf, err := json.Marshal(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(buf)
	})

	http.HandleFunc("/key/{name}", func(w http.ResponseWriter, req *http.Request) {
		v, ok := cache.Load(req.PathValue("name"))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		buf, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(buf)
	})

	log.Printf("starting listening http server on %s", cfg.Http.Addr)
	if err := http.ListenAndServe(cfg.Http.Addr, nil); err != nil {
		log.Fatalf("error listing on '%s': %v", cfg.Http.Addr, err)
	}
}

func scrapeMetrics(cfg *pkg.Config, cache *pkg.Cache, a v1.API) {
	for _, expr := range cfg.Expressions {
		res, warnings, err := a.Query(context.Background(), expr.Query, time.Now())
		if err != nil {
			log.Printf("error running query '%s': %v", expr.Query, err)
			break
		}

		if len(warnings) > 0 {
			log.Printf("warnings running query '%s': %v", expr.Query, warnings)
		}

		out, err := executeExpression(res, expr.Experssion)
		if err != nil {
			log.Printf("error executing query '%s' with expression '%s': %v", expr.Query, expr.Experssion, err)
			continue
		}

		cache.Store(expr.Name, out)
	}
}

func executeExpression(res model.Value, exp string) (any, error) {
	if exp == "" {
		return res, nil
	}

	env := map[string]interface{}{
		"query_result": res,
	}

	program, err := expr.Compile(exp, expr.Env(env))
	if err != nil {
		return nil, err
	}

	out, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}

	return out, nil
}
