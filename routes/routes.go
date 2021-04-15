package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	// 微信
	rg := engine.Group("/api")

	// publish
	RegPublish(rg)
}
