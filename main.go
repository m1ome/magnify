package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/m1ome/magnify/pkg"
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

	log.Print("scraping metrics on a startup")
	cache := pkg.NewCache(time.Second * time.Duration(cfg.Cache.Expiration))
	metrics, err := pkg.NewMetrics(cfg, cache)
	if err != nil {
		log.Fatalf("error initializing metrics: %v", err)
	}

	runUpdate(metrics)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		sendResponse(cache.Copy(), w)
	})

	http.HandleFunc("/{name}", func(w http.ResponseWriter, req *http.Request) {
		v, ok := cache.Load(req.PathValue("name"))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		sendResponse(v, w)
	})

	log.Printf("starting listening http server on %s", cfg.Http.Addr)
	if err := http.ListenAndServe(cfg.Http.Addr, nil); err != nil {
		log.Fatalf("error listing on '%s': %v", cfg.Http.Addr, err)
	}
}

func sendResponse(v any, w http.ResponseWriter) {
	buf, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(buf)
}

func runUpdate(metrics *pkg.Metrics) {
	if err := metrics.Update(); err != nil {
		log.Fatalf("error updating metrics: %v", err)
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			if err := metrics.Update(); err != nil {
				log.Printf("error updating metrics: %v", err)
			}
		}
	}()
}
