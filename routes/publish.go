package routes

import (
	"github.com/gin-gonic/gin"
)

func RegPublish(rg *gin.RouterGroup) {

	// 市场状态
	//rg.POST("/market", publish.NewMarketStatus().POST) //默认pb模式

	// 移动端首页
	//rg.GET("/mindex", publish.NewMIndex().GET)
}
