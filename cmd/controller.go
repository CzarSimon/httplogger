package main

import (
	"net/http"

	"github.com/CzarSimon/httplogger/pkg/httputil"
	"github.com/CzarSimon/httplogger/pkg/log"
	"github.com/CzarSimon/httplogger/pkg/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func handleLog(c *gin.Context) {
	event, err := getEvent(c)
	if err != nil {
		c.Error(err)
		return
	}

	go log.Log(event)
	httputil.SendOK(c)
}

func getEvent(c *gin.Context) (*models.Event, error) {
	var event models.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		logger.Error("Failed to parse log event", zap.Error(err))
		return nil, httputil.NewError("Failed to parse log event", http.StatusBadRequest)
	}
	return &event, nil
}
