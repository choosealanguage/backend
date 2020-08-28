package webserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func fail(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"code":  code,
		"error": message,
	})
}

func failNotFound(ctx *gin.Context) {
	fail(ctx, http.StatusNotFound, "not found")
}

func failInternal(ctx *gin.Context, err error) {
	fail(ctx, http.StatusInternalServerError, err.Error())
}

func failBadRequest(ctx *gin.Context) {
	fail(ctx, http.StatusBadRequest, "bad request")
}

func ok(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "ok",
	})
}

func (ws *WebServer) handlerCORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", ws.config.CorsOrigin)
	ctx.Header("Access-Control-Allow-Methods", "GET")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,Cookie")
	ctx.Header("Access-Control-Allow-Credentials", "true")
}

func (ws *WebServer) handlePreflight(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusOK)
		ctx.Abort()
		return
	}
	ctx.Next()
}
