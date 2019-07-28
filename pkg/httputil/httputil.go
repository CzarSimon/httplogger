package httputil

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const metricsPath = "/metrics"

var errLog = GetLogger("errorLog")

// Prometheus metrics.
var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number served requests",
		},
		[]string{"endpoint", "method", "status"},
	)
	requestsLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_latency_ms",
			Help: "Request latency in milliseconds",
		},
		[]string{"endpoint", "method", "status"},
	)
)

// Error implements the error interface with a message and http status code.
type Error struct {
	ID         string `json:"id,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status,omitempty"`
}

// NewError creates a error with a message and a status code.
func NewError(message string, status int) *Error {
	if message == "" {
		message = http.StatusText(status)
	}

	return &Error{
		ID:         newID(),
		Message:    message,
		StatusCode: status,
	}
}

// BadRequest creates a bad request error.
func BadRequest(message string) *Error {
	return NewError(message, http.StatusBadRequest)
}

// InternalServerError creates an internal server error.
func InternalServerError(message string) *Error {
	return NewError(message, http.StatusInternalServerError)
}

// SendOK sends an ok status and message to the client.
func SendOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

// SendError formats, logs and sends a response back to the client
func SendError(err *Error, c *gin.Context) {
	c.AbortWithStatusJSON(err.StatusCode, err)
}

// NewRouter creates a default router.
func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), HandleErrors(), Metrics())
	r.GET("/health", SendOK)
	r.GET("/metrics", prometheusHandler())

	return r
}

// HandleErrors wrapper function to deal with encountered errors
// during request handling.
func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := getFirstError(c)
		if err == nil {
			return
		}

		var httpError *Error
		switch err.(type) {
		case *Error:
			httpError = err.(*Error)
			break
		default:
			httpError = InternalServerError(err.Error())
			break
		}

		logError(httpError)
		SendError(httpError, c)
	}
}

// Metrics records metrics about a request.
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == metricsPath {
			c.Next()
			return
		}
		stop := createTimer()
		endpoint := c.FullPath()
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		latency := stop()
		requestsTotal.WithLabelValues(endpoint, method, status).Inc()
		requestsLatency.WithLabelValues(endpoint, method, status).Observe(latency)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// logError logs internal errors.
func logError(err *Error) {
	if err.StatusCode < 500 {
		errLog.Info(err.Message, zap.Int("status", err.StatusCode), zap.String("errorId", err.ID))
		return
	}

	errLog.Error(err.Message, zap.Int("status", err.StatusCode), zap.String("errorId", err.ID))
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error(id=[%s], statusCode=%d, message=%s)", e.ID, e.StatusCode, e.Message)
}

// GetLogger creates a named logger for internal application logs.
func GetLogger(name string) *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln("Failed to get zap.Logger", err)
	}
	return logger.With(zap.String("app", "httplogger"), zap.String("logger", name))
}

// getFirstError returns the first error in the gin.Context, nil if not present.
func getFirstError(c *gin.Context) error {
	allErrors := c.Errors
	if len(allErrors) == 0 {
		return nil
	}
	return allErrors[0].Err
}

type calcDuration func() float64

func createTimer() calcDuration {
	start := time.Now()

	// Returns latency in milliseconds.
	return func() float64 {
		end := time.Now()
		return float64(end.Sub(start)) / 1e6
	}
}

func newID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}

	return id.String()
}
