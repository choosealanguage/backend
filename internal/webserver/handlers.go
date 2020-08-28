package webserver

import (
	"net/http"

	"github.com/choosealanguage/backend/pkg/provider"
	"github.com/gin-gonic/gin"
)

func (ws *WebServer) handlerGetLanguages(ctx *gin.Context) {
	langMap := ws.prov.GetLanguages()
	langList := make([]*provider.LanguageModel, len(langMap))

	i := 0
	for _, v := range langMap {
		langList[i] = v
		i++
	}

	ctx.JSON(http.StatusOK, langList)
}

func (ws *WebServer) handlerGetLanguage(ctx *gin.Context) {
	id := ctx.Param("id")

	lang, ok := ws.prov.GetLanguages()[id]
	if !ok {
		failNotFound(ctx)
		return
	}

	ctx.JSON(http.StatusOK, lang)
}
