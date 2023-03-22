package api

import (
	"subscribe2clash/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"subscribe2clash/internal/clash"
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
	sub, err := model.GetSubscribeByShortCode(linksName)
	if err != nil {
		c.String(http.StatusBadRequest, key+"值为空")
		c.Abort()
		return
	}

	config, err := clash.Config(clash.Url, sub.SubscribeURL)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}

func (cc *ClashController) Nodes(c *gin.Context) {

	nodes, err := model.GetAllNodes()
	if err != nil {
		c.String(http.StatusBadRequest, key+"值为空")
		c.Abort()
		return
	}

	urls := make([]string, 0)
	for _, node := range nodes {
		urls = append(urls, node.Address)
	}
	config, err := clash.Nodes(urls)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	c.String(http.StatusOK, config)
}
