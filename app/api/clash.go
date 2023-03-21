package api

import (
	"github.com/icpd/subscribe2clash/internal/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpd/subscribe2clash/internal/clash"
)

const key = "link"

type ClashController struct{}

func (cc *ClashController) Clash(c *gin.Context) {
	links, exists := c.GetQuery(key)
	if !exists {
		links, _ = c.GetQuery("sub_link") // 兼容旧key
	}

	if links == "" {
		c.String(http.StatusBadRequest, key+"不能为空")
		c.Abort()
		return
	}

	config, err := clash.Config(clash.Url, links)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}

func (cc *ClashController) Self(c *gin.Context) {
	linksName, exists := c.GetQuery(key)
	if !exists {
		linksName, _ = c.GetQuery("sub_link") // 兼容旧key
	}

	if linksName == "" {
		c.String(http.StatusBadRequest, key+"不能为空")
		c.Abort()
		return
	}
	links := global.Subscribes[linksName]
	if linksName == "" {
		c.String(http.StatusBadRequest, key+"值为空")
		c.Abort()
		return
	}
	config, err := clash.Config(clash.Url, links)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}
