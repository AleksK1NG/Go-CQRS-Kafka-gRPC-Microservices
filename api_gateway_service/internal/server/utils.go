package server

import (
	"context"
	"github.com/heptiolabs/healthcheck"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddReadinessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		if s.cfg != nil {
			return nil
		}
		return errors.New("Config not loaded")
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	go func() {
		s.log.Infof("API Gateway Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}

func (s *server) runMetrics(cancel context.CancelFunc) {
	metricsServer := echo.New()
	metricsServer.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         stackSize,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	go func() {
		metricsServer.GET(s.cfg.Probes.PrometheusPath, echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Probes.PrometheusPort)
		if err := metricsServer.Start(s.cfg.Probes.PrometheusPort); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()
}
