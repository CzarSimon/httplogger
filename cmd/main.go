package main

import (
	"fmt"
	"net/http"

	"github.com/CzarSimon/httputil"
	logutil "github.com/CzarSimon/httputil/logger"
	"go.uber.org/zap"
)

var logger = logutil.GetDefaultLogger("httplogger/main")

const port = ":8080"

func main() {
	logger.Info(fmt.Sprintf("Starting httplogger on port: %s", port))
	err := server().ListenAndServe()
	if err != nil {
		logger.Error("Unexpected server error", zap.Error(err))
	}
}

func server() *http.Server {
	r := httputil.NewRouter("httplogger", func() error {
		return nil
	})
	r.Use(httputil.AllowJSON())

	r.POST("/v1/logs", handleLog)

	return &http.Server{
		Addr:    port,
		Handler: r,
	}
}
