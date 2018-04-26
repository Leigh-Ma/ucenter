package api

import (
	"ucenter/controllers"
)

type itemController struct {
	authorizedController
}

func (c *itemController) Export() func(string) {
	return controllers.Export(c, map[string]string{
	})
}
