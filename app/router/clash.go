package router

import (
	"github.com/gin-gonic/gin"

	"subscribe2clash/app/api"
)

func clashRouter(r *gin.Engine) {
	clash := api.ClashController{}
	r.GET("/self", clash.Self)
	r.GET("/nodes", clash.Nodes)

	r.GET("/", clash.Clash)

}
