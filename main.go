package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	fmt.Printf("%v | %v | %v | %v | %v | %v \n",
		timeIn.Format(time.RFC3339),
		c.RealIP(),
		c.Response().Status,
		c.Request().Method,
		c.Request().URL.Path,
		currentTime.Sub(timeIn).String(),
	)
}

type (
	Stats struct {
		Uptime       time.Time      `json:"uptime_since"`
		RequestCount uint64         `json:"request_count"`
		Statuses     map[string]int `json:"statuses"`
		IPAddresses  map[string]int `json:"requests_by_ip_address"`
		mutex        sync.RWMutex
	}

	ImgInfo struct {
		Path    string
		Caption string
	}
)

func main() {
	e := echo.New()
	s := NewStats()

	masterTemplate, _ = template.ParseFiles("assets/template.html")

	e.Use(s.Process)
	e.Use(middleware.Recover())

	//setImg(e)
	setRoutes(e, s)

	e.Logger.Fatal(e.Start(":3000"))
}

func setRoutes(e *echo.Echo, s *Stats) {

	e.GET("/resume", func(c echo.Context) error {
		return c.File("assets/mannes_resume.pdf")
	})

	e.GET("/n", func(c echo.Context) error {
		return c.File("assets/n.png")
	})

	e.GET("/m", func(c echo.Context) error {
		return c.File("assets/m.png")
	})

	e.GET("/style.css", func(c echo.Context) error {
		return c.File("assets/style.css")
	})

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, s, "   ")
	})

	e.GET("/*", Route)

}

func setImg(e *echo.Echo) error {
	root := "assets/img"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".jpg") {
			e.GET(path, func(c echo.Context) error {
				return c.File(path)
			})

			sections := strings.Split(path, "/")
			fileName := strings.Split(sections[2], ".")
			Images = append(Images, ImgInfo{
				Path:    fmt.Sprintf("\"%v\"", path),
				Caption: strings.ReplaceAll(fileName[0], "_", " "),
			})

		}
		return nil
	})
	if err != nil {
	}
	return nil
}
