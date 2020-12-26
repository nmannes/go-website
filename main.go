package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
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

	masterTemplate, _ = template.ParseFiles("assets/template.html")

	setRoutes(e)

	e.Logger.Fatal(e.Start(":8000"))
}

func serveFileWithCache(e *echo.Echo, pathToFile, route string) error {
	f, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return err
	}

	e.GET(route, func(c echo.Context) error {
		return c.Blob(
			http.StatusOK,
			http.DetectContentType(f),
			f,
		)
	})

	return nil
}

func setRoutes(e *echo.Echo) {
	s := NewStats()

	e.Use(s.Process)
	e.Use(middleware.Recover())

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, s, "\t")
	})

	serveFileWithCache(e, "assets/style.css", "/style")
	serveFileWithCache(e, "assets/mannes_resume.pdf", "/resume")

	setIcons(e)
	setImg(e)

	e.GET("/*", Route)

}

func setIcons(e *echo.Echo) error {

	mFile, err := ioutil.ReadFile("assets/m.png")
	if err != nil {
		return err
	}

	nFile, err := ioutil.ReadFile("assets/n.png")
	if err != nil {
		return err
	}

	e.GET("/favicon.ico", func(c echo.Context) error {
		fileReturn := mFile
		if rand.Intn(2) == 0 {
			fileReturn = nFile
		}

		return c.Blob(
			http.StatusOK,
			http.DetectContentType(fileReturn),
			fileReturn,
		)
	})

	return nil
}

func setImg(e *echo.Echo) error {
	root := "assets/img"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(info.Name(), ".jpg") ||
			strings.Contains(info.Name(), ".png") {

			sections := strings.Split(path, "/")
			fileName := strings.Split(sections[2], ".")
			imageInfo := ImgInfo{
				Path:    fmt.Sprintf("\"%v\"", path),
				Caption: strings.ReplaceAll(fileName[0], "_", " "),
			}
			Images = append(Images, imageInfo)

			serveFileWithCache(e, path, path)

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
