package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now().UTC(),
		Statuses: map[string]int{},
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		return nil
	}
}

type (
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

func main() {
	// Echo instance
	s := NewStats()
	e := echo.New()

	// Middleware
	e.Use(s.Process)
	e.Use(middleware.Recover())

	e.GET("/resume", func(c echo.Context) error {
		return c.Inline("mannes_resume.pdf", "Nathan's resume")
	})

	e.GET("/nm", func(c echo.Context) error {
		return c.File("nm.gif")
	})

	e.GET("/n", func(c echo.Context) error {
		return c.File("n.png")
	})
	e.GET("/m", func(c echo.Context) error {
		return c.File("m.png")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, s, "   ")
	})

	e.GET("/*", func(c echo.Context) error {
		return c.File("index.html")
	})

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}
