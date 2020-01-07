package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(ctx *gin.Context)  {
	load := ctx.Query("load")
	ctx.String(http.StatusOK, "welcome " + load)
}