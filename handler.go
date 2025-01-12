package elements_echo

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	htmlt "html/template"
	"net/http"
	"os"
	"strings"
)

//go:embed template/index.html
var template string

//go:embed template/favicon.png
var icon string

//go:embed template/styles.min.css
var styles string

//go:embed template/web-components.min.js
var script string

type StopLightMiddleware struct {
	urlPrefix string
	specRoot  string
}

func New(urlPrefix, specRoot string) *StopLightMiddleware {
	return &StopLightMiddleware{urlPrefix: urlPrefix, specRoot: specRoot}
}

func (m *StopLightMiddleware) Handle() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			m.httpHandler(ctx.Response(), ctx.Request())
			if ctx.Response().Committed {
				return nil
			}
			return next(ctx)
		}
	}
}

func (m *StopLightMiddleware) template() ([]byte, error) {
	var t, err = htmlt.New("index.html").Parse(template)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, map[string]string{"Path": m.urlPrefix})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *StopLightMiddleware) getSpec() ([]byte, error) {
	specFile := fmt.Sprintf("%s%s%s", m.specRoot, string(os.PathSeparator), "swagger.yaml")
	if specFile == "" {
		panic(errors.New("spec file not exist"))
	}
	var spec []byte
	spec, err := os.ReadFile(specFile)
	if err != nil {
		return nil, err
	}
	return spec, nil
}

func (m *StopLightMiddleware) httpHandler(w http.ResponseWriter, req *http.Request) {
	spec, err := m.getSpec()
	if err != nil {
		return
	}
	data, err := m.template()
	if err != nil {
		return
	}
	
	method := strings.ToLower(req.Method)
	if method != "get" && method != "head" {
		return
	}
	
	header := w.Header()
	
	strippedPrefix := strings.TrimPrefix(req.URL.Path, m.urlPrefix)
	if strippedPrefix == "/" || strippedPrefix == "/index.html" || strippedPrefix == "" {
		header.Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
		return
	}
	if strippedPrefix == "/swagger.yaml" {
		header.Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(spec)
		return
	}
	if strippedPrefix == "/script.js" {
		header.Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(script))
		return
	}
	if strippedPrefix == "/styles.css" {
		header.Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(styles))
		return
	}
	if strippedPrefix == "/favicon.png" {
		header.Set("Content-Type", "image/png")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(icon))
		return
	}
	
}
