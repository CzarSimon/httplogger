package main

import (
	"fmt"
	"net/http"

	"github.com/CzarSimon/httplogger/pkg/httputil"
	"go.uber.org/zap"
)

var logger = httputil.GetLogger("main")

const port = ":8080"

func main() {
	logger.Info(fmt.Sprintf("Starting httplogger on port: %s", port))
	err := server().ListenAndServe()
	if err != nil {
		logger.Error("Unexpected server error", zap.Error("err", err))
	}
}

func server() *http.Server {
	r := httputil.NewRouter()
	r.POST("/v1/logs", handleLog)
	r.GET("/health", httputil.SendOK)

	return &http.Server{
		Addr:    port,
		Handler: r,
	}
}
