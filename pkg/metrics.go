package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/expr-lang/expr"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type (
	Metrics struct {
		c     *Config
		api   v1.API
		cache *Cache
	}
)

func NewMetrics(config *Config, cache *Cache) (*Metrics, error) {
	c, err := api.NewClient(api.Config{
		Address: config.Prometheus.Addr,
	})
	if err != nil {
		return nil, fmt.Errorf("error connecting to prometheus: %v", err)
	}

	return &Metrics{
		c:     config,
		api:   v1.NewAPI(c),
		cache: cache,
	}, nil
}

func (m *Metrics) Update() error {
	if err := m.scrapeMetrics(); err != nil {
		return err
	}

	m.cache.Cleanup()
	return nil
}

func (m *Metrics) scrapeMetrics() error {
	for _, expr := range m.c.Expressions {
		res, _, err := m.api.Query(context.Background(), expr.Query, time.Now())
		if err != nil {
			return err
		}

		out, err := m.executeExpression(res, expr.Experssion)
		if err != nil {
			return err
		}

		m.cache.Store(expr.Name, out)
	}

	return nil
}

func (m *Metrics) executeExpression(res model.Value, exp string) (any, error) {
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
