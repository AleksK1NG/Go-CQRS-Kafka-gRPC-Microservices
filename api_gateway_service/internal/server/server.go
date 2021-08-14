package server

import (
	"context"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/config"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/client"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/metrics"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/middlewares"
	v1 "github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/products/delivery/http/v1"
	"github.com/AleksK1NG/cqrs-microservices/api_gateway_service/internal/products/service"
	"github.com/AleksK1NG/cqrs-microservices/pkg/interceptors"
	"github.com/AleksK1NG/cqrs-microservices/pkg/kafka"
	"github.com/AleksK1NG/cqrs-microservices/pkg/logger"
	"github.com/AleksK1NG/cqrs-microservices/pkg/tracing"
	readerService "github.com/AleksK1NG/cqrs-microservices/reader_service/proto/product_reader"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log  logger.Logger
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	im   interceptors.InterceptorManager
	echo *echo.Echo
	ps   *service.ProductService
	m    *metrics.ApiGatewayMetrics
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, echo: echo.New(), v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg)
	s.im = interceptors.NewInterceptorManager(s.log)
	s.m = metrics.NewApiGatewayMetrics(s.cfg)

	readerServiceConn, err := client.NewReaderServiceConn(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer readerServiceConn.Close() // nolint: errcheck
	rsClient := readerService.NewReaderServiceClient(readerServiceConn)

	kafkaProducer := kafka.NewProducer(s.log, []string{"host.docker.internal:9092"})
	defer kafkaProducer.Close() // nolint: errcheck

	s.ps = service.NewProductService(s.log, s.cfg, kafkaProducer, rsClient)

	productHandlers := v1.NewProductsHandlers(s.echo.Group(s.cfg.Http.ProductsPath), s.log, s.mw, s.cfg, s.ps, s.v, s.m)
	productHandlers.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("API Gateway is listening on PORT: %s", s.cfg.Http.Port)

	s.runMetrics(cancel)
	s.runHealthCheck(ctx)

	if s.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
		if err != nil {
			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}

	s.log.Info("API Gateway server exited properly")
	return nil
}
