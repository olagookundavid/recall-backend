package api

import (
	"bufio"
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"recall-app/internal/token"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tomasen/realip"
	"tailscale.com/tstime/rate"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func (app *Application) RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// Use the recover function to check if there has been a panic or not.
			if err := recover(); err != nil {
				c.Writer.Header().Set("Connection", "close")
				app.ServerErrorResponse(c, fmt.Errorf("%s", err))
			}
		}()
		c.Next()
	}
}

func (app *Application) RateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background cleanup goroutine
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		if !app.Config.Limiter.Enabled {
			c.Next()
			return
		}

		ip := realip.FromRequest(c.Request)
		mu.Lock()

		// Initialize or get existing limiter for this IP
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Limit(app.Config.Limiter.Rps),
					app.Config.Limiter.Burst),
			}
		}
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(c)
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}

// AuthMiddleware creates a gin middleware for authorization
func (app *Application) TokenMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			app.unAuthorizedResponse(c, err.Error())
			c.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			app.unAuthorizedResponse(c, err.Error())
			c.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			app.unAuthorizedResponse(c, err.Error())
			c.Abort()
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			app.unAuthorizedResponse(c, err.Error())
			c.Abort()
			return
		}
		c.Set(authorizationPayloadKey, payload)
		c.Next()
	}
}

func (app *Application) Metrics() gin.HandlerFunc {
	var (
		totalRequestsReceived           = expvar.NewInt("total_requests_received")
		totalResponsesSent              = expvar.NewInt("total_responses_sent")
		totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
		totalResponsesSentByStatus      = expvar.NewMap("total_responses_sent_by_status")
	)
	// The following code will be run for every request...
	return func(c *gin.Context) {
		// Record the time that we started to process the request.
		start := time.Now()
		totalRequestsReceived.Add(1)
		mw := &metricsResponseWriter{ResponseWriter: c.Writer}
		c.Writer = mw

		c.Next()

		totalResponsesSent.Add(1)
		totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	}

}

type metricsResponseWriter struct {
	gin.ResponseWriter
	statusCode    int
	headerWritten bool
}

// Override `WriteHeader` to track status code
func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
	if !mw.headerWritten {
		mw.statusCode = statusCode
		mw.headerWritten = true
	}
	mw.ResponseWriter.WriteHeader(statusCode)
}

// Override `Write` to track response body writes
func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
	if !mw.headerWritten {
		mw.statusCode = http.StatusOK
		mw.headerWritten = true
	}
	return mw.ResponseWriter.Write(b)
}

// Implement `Unwrap()` to return the underlying writer
func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
	return mw.ResponseWriter
}

// Implement `Flush()`, `Hijack()`, and `CloseNotify()` for full compatibility

func (mw *metricsResponseWriter) Flush() {
	if flusher, ok := mw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (mw *metricsResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := mw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

func (mw *metricsResponseWriter) CloseNotify() <-chan bool {
	if notifier, ok := mw.ResponseWriter.(http.CloseNotifier); ok {
		return notifier.CloseNotify()
	}
	return nil
}

func (app *Application) WrapHTTPHandler(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Convert gin.Context to http.Request and ResponseWriter
		h.ServeHTTP(c.Writer, c.Request)
	}
}
