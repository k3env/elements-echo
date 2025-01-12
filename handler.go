package elements_echo

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/labstack/echo/v4"
	htmlt "html/template"
	"net/http"
	"os"
	"path/filepath"
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
	urlPrefix   string
	specContent []byte
	specFormat  string
}

func (m *StopLightMiddleware) UseSpecFile(path string) (*StopLightMiddleware, error) {
	ext := filepath.Ext(path)
	if ext == ".json" {
		m.specFormat = "json"
	} else if ext == ".yaml" || ext == ".yml" {
		m.specFormat = "yaml"
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m.specContent = content

	return m, nil
}

func (m *StopLightMiddleware) UseEmbed(fs embed.FS, path string) (*StopLightMiddleware, error) {
	ext := filepath.Ext(path)
	if ext == ".json" {
		m.specFormat = "json"
	} else if ext == ".yaml" || ext == ".yml" {
		m.specFormat = "yaml"
	}

	content, err := fs.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m.specContent = content
	return m, nil
}

func (m *StopLightMiddleware) UseContent(content []byte, format string) *StopLightMiddleware {
	m.specFormat = format
	m.specContent = content
	return m
}

func New(urlPrefix string) *StopLightMiddleware {
	return &StopLightMiddleware{urlPrefix: urlPrefix}
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

func (m *StopLightMiddleware) httpHandler(w http.ResponseWriter, req *http.Request) {
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
	if strippedPrefix == "/swagger.yaml" && m.specFormat == "yaml" {
		header.Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(m.specContent)
		return
	}
	if strippedPrefix == "/swagger.yml" && m.specFormat == "yaml" {
		header.Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(m.specContent)
		return
	}
	if strippedPrefix == "/swagger.json" && m.specFormat == "json" {
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(m.specContent)
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
