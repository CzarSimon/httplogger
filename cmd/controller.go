package main

import (
	"fmt"

	"github.com/CzarSimon/httplogger/internal/log"
	"github.com/CzarSimon/httplogger/internal/models"
	"github.com/CzarSimon/httputil"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func handleLog(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "controller_handle_log")
	defer span.Finish()

	event, err := getEvent(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = validateEvent(event)
	if err != nil {
		c.Error(err)
		return
	}

	go log.Log(ctx, event)
	httputil.SendOK(c)
}

func getEvent(c *gin.Context) (*models.Event, error) {
	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		err := fmt.Errorf("failed to parse log event: %w", err)
		return nil, httputil.BadRequestError(err)
	}
	return &event, nil
}

func validateEvent(e *models.Event) error {
	if e.Level == models.DebugLevel || e.Level == models.InfoLevel || e.Level == models.WarnLevel || e.Level == models.ErrorLevel {
		return nil
	}

	err := fmt.Errorf("unsupported log level: %s", e.Level)
	return httputil.BadRequestError(err)
}
