package webserver

import (
	"github.com/choosealanguage/backend/pkg/provider"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Debug      bool
	Address    string
	CorsOrigin string
}

type WebServer struct {
	router *gin.Engine
	prov   *provider.Provider
	config *Config
}

func New(prov *provider.Provider, config Config) (ws *WebServer) {
	ws = new(WebServer)

	ws.prov = prov
	ws.config = &config
	ws.router = gin.Default()

	ws.registerHandlers()

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return
}

func (ws *WebServer) Run() error {
	return ws.router.Run(ws.config.Address)
}

func (ws *WebServer) registerHandlers() {
	if ws.config.Debug {
		ws.router.Use(ws.handlerCORS)
	}

	ws.router.Use(ws.handlePreflight)

	langs := ws.router.Group("/languages")
	langs.
		GET("", ws.handlerGetLanguages).
		GET("/:id", ws.handlerGetLanguage)
}
