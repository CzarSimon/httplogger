package httputil

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var errLog = GetLogger("errorLogger")

// Error implements the error interface with a message and http status code.
type Error struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"status,omitempty"`
}

// NewError creates a error with a message and a status code.
func NewError(messge string, status int) *Error {
	return &Error{
		Message:    messge,
		StatusCode: status,
	}
}

// NewInternalError creates an internal server error.
func NewInternalError(message string) *Error {
	return &Error{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
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
	r.Use(gin.Recovery(), HandleErrors())

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
			httpError = NewInternalError(err.Error())
			break
		}

		logError(httpError)
		SendError(httpError, c)
	}
}

// logError logs internal errors.
func logError(err *Error) {
	if err.StatusCode < 500 {
		errLog.Info(err.Message, zap.Int("status", err.StatusCode))
		return
	}

	errLog.Error(err.Message, zap.Int("status", err.StatusCode))
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error(statusCode=%d, message=%s)", e.StatusCode, e.Message)
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
