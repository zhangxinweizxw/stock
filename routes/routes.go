package routes

import (
	"github.com/gin-gonic/gin"
)

func Register(engine *gin.Engine) {

	rg := engine.Group("/api")

	// publish
	RegPublish(rg)
}
