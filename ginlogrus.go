// Package ginlogrus provides a logging middleware to get
// https://github.com/sirupsen/logrus as logging library for
// https://github.com/gin-gonic/gin. It can be used as replacement for
// the internal logging middleware
// http://godoc.org/github.com/gin-gonic/gin#Logger.
//
// Derived on https://github.com/zalando/gin-glog
//
// Example:
//    package main
//    import (
//        "flag"
//        "time"
//        log "github.com/sirupsen/logrus"
//        "github.com/rocksolidlabs/gin-logrus"
//        "github.com/gin-gonic/gin"
//    )
//    func main() {
//        flag.Parse()
//        router := gin.New()
//        router.Use(ginlogrus. Logger("MYAPI", false, true, os.Stdout, logrus.WarnLevel))
//        //..
//        router.Use(gin.Recovery())
//        log.Info("bootstrapped application")
//        router.Run(":8080")
//    }
//
package ginlogrus

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"

	"gopkg.in/gin-gonic/gin.v1"
)

var log *logrus.Logger

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// ErrorLogger returns an ErrorLoggerT with parameter gin.ErrorTypeAny
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns an ErrorLoggerT middleware with the given
// type gin.ErrorType.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

// Logger prints a logline for each request and measures the time to
// process for a call. It formats the log entries similar to
// http://godoc.org/github.com/gin-gonic/gin#Logger does.
//
// Example:
//        router := gin.New()
//        router.Use(ginlogrus.Logger(false, true, os.Stdout, log.WarnLevel))
func Logger(l *logrus.Logger, outputTag string, outputJSON bool, outputColor bool, outputFile io.Writer, outLevel logrus.Level) gin.HandlerFunc {

	// set the logger
	log = l

	// Set the output tag
	if outputTag == "" {
		outputTag = "GIN"
	}

	// Log as JSON instead of the default ASCII formatter.
	if outputJSON {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		reset = ""
	}

	// Turn off logrus color
	if !outputColor && !outputJSON {
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, DisableColors: true})
	}

	// Output to stdout instead of the default stderr, could also be a file.
	logrus.SetOutput(outputFile)

	// Set log severity oputLevel or above.
	logrus.SetLevel(outLevel)

	return func(c *gin.Context) {
		t := time.Now()

		// process request
		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		statusColor := reset
		methodColor := reset
		if outputColor {
			statusColor = colorForStatus(statusCode)
			methodColor = colorForMethod(method)
		}
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				log.Warningf("[%s] |%s %3d %s| %12v | %s |%s  %s %-7s %s %s",
					outputTag,
					statusColor, statusCode, reset,
					latency,
					clientIP,
					methodColor, reset, method,
					path,
					c.Errors.String(),
				)
			}
		case statusCode >= 500:
			{
				log.Errorf("[%s] |%s %3d %s| %12v | %s |%s  %s %-7s %s %s",
					outputTag,
					statusColor, statusCode, reset,
					latency,
					clientIP,
					methodColor, reset, method,
					path,
					c.Errors.String(),
				)
			}
		default:
			log.Debugf("[%s] |%s %3d %s| %12v | %s |%s  %s %-7s %s\n%s",
				outputTag,
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, reset, method,
				path,
				c.Errors.String(),
			)
		}

	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}
