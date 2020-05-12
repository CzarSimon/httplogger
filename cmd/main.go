package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/CzarSimon/httputil"
	logutil "github.com/CzarSimon/httputil/logger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var logger = logutil.GetDefaultLogger("httplogger/main")

const port = ":8080"

func main() {
	traceCloser := setupTracer()
	defer traceCloser.Close()

	s := server()
	logger.Info(fmt.Sprintf("Starting httplogger on port: %s", port))

	err := s.ListenAndServe()
	if err != nil {
		logger.Error("Unexpected server error", zap.Error(err))
	}
}

func server() *http.Server {
	r := httputil.NewCustomRouter(
		healthCheck,
		gin.Recovery(),
		httputil.Trace("httplogger"),
		httputil.Metrics(),
		httputil.HandleErrors(),
	)
	r.Use(httputil.AllowJSON())

	r.POST("/v1/logs", handleLog)

	return &http.Server{
		Addr:    port,
		Handler: r,
	}
}

func setupTracer() io.Closer {
	jcfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Fatal("failed to create jaeger configuration", zap.Error(err))
	}

	tracer, closer, err := jcfg.NewTracer()
	if err != nil {
		log.Fatal("failed to create tracer", zap.Error(err))
	}

	opentracing.SetGlobalTracer(tracer)
	return closer
}

func healthCheck() error {
	return nil
}
