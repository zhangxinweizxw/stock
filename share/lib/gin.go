package lib

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Json(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func WriteString(gc *gin.Context, code int, data interface{}) {
	res := map[string]interface{}{"code": code}
	if data != nil {
		res["data"] = data
	}
	gc.JSON(http.StatusOK, res)
}

func WriteData(gc *gin.Context, data []byte) {

	gc.Data(http.StatusOK, "application/octet-stream", data)
}
