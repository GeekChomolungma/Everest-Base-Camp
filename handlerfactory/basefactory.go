package handlerfactory

import (
	"net/http"

	"github.com/GeekChomolungma/Everest-Base-Camp/dtos"
	"github.com/gin-gonic/gin"
)

type BaseFactory interface {
	Create() BaseHandler
}

type BaseHandler interface {
	Get(url string) (string, error)
	Post(url string, body string) (string, error)
}

type ImportIns struct {
	AimSite string `json:"aimsite"`
	Method  string `json:"method"`
	Url     string `json:"url"`
	Body    string `json:"body"`
}

func FactoryImport(c *gin.Context) {
	var importReq ImportIns
	err := c.Bind(&importReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": dtos.CANNOT_PARSE_POST_BODY, "msg": "Sorry", "data": err.Error()})
		return
	}

	if factory, ok := FactoryInstanceMap[importReq.AimSite]; ok {
		handler := factory.Create()
		switch importReq.Method {
		case "GET":
			result, err := handler.Get(importReq.Url)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.AIM_SITE_GET_ERROR, "msg": "wrong", "data": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.OK, "msg": "succeed", "data": result})
			}
		case "POST":
			result, err := handler.Post(importReq.Url, importReq.Body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.AIM_SITE_POST_ERROR, "msg": "wrong", "data": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"code": dtos.OK, "msg": "succeed", "data": result})
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"code": dtos.UNKNOW_METHOD, "msg": "Sorry", "data": ""})
		}
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"code": dtos.AIM_SITE_NOT_EXIST, "msg": "Sorry", "data": err.Error()})
}
