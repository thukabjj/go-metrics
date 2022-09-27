package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}
