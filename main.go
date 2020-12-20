package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"go.uber.org/zap"
)

func NewStats() *Stats {
	return &Stats{
		Uptime:      time.Now().UTC(),
		Statuses:    map[string]int{},
		IPAddresses: map[string]int{},
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		timeIn := time.Now().UTC()
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		s.IPAddresses[c.RealIP()]++

		log(c, timeIn, time.Now().UTC())
		return nil
	}
}

func log(c echo.Context, timeIn time.Time, currentTime time.Time) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logger.Info("new request",
		zap.String("requested_at", timeIn.Format(time.RFC3339)),
		zap.String("ip", c.RealIP()),
		zap.Int("status", c.Response().Status),
		zap.String("path", c.Request().URL.Path),
		zap.String("response_time", currentTime.Sub(timeIn).String()),
	)
}

type (
	Stats struct {
		Uptime       time.Time      `json:"uptime_since"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		IPAddresses  map[string]int `json:"requests_by_ip_address"`
		mutex        sync.RWMutex
	}
)

func main() {
	s := NewStats()
	e := echo.New()

	e.Use(s.Process)
	e.Use(middleware.Recover())

	files := []string{
		"mannes_resume.pdf",
		"m.png",
		"n.png",
	}

	for _, f := range files {
		e.GET(fmt.Sprintf("/%v", f), func(c echo.Context) error {
			return c.File(f)
		})
	}

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, s, "   ")
	})

	e.GET("/*", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.Logger.Fatal(e.Start(":3000"))
}
