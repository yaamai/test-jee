package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Handler struct {
	DataRoot string
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := Handler{DataRoot: "data"}
	e.GET("/api/data", handler.GetData)
	e.GET("/*", echo.NotFoundHandler, middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "static",
		Browse: false,
		HTML5:  true,
	}))
	e.Logger.Fatal(e.Start(":18328"))
}

func (h *Handler) GetData(c echo.Context) error {
	// currently, expect fixed target json files structure
	files, err := os.ReadDir(h.DataRoot)
	if err != nil {
		return err
	}

	resp := []interface{}{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		path := filepath.Join(h.DataRoot, file.Name(), "data.json")

		dataBytes, err := os.ReadFile(path)
		if err != nil {
			log.Println(err)
			continue
		}

		data := map[string]interface{}{}
		err = json.Unmarshal(dataBytes, &data)
		if err != nil {
			log.Println(err)
			continue
		}

		resp = append(resp, data)
	}

	return c.JSON(http.StatusOK, resp)
}
