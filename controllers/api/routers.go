package api

import "ucenter/controllers"

var ApiRouter = &controllers.RouterGroup{
	Namespace: "/api/authed/",
	Routers: map[string]controllers.IExport{
		"sign":   &signController{},
		"player": &playerController{},
		"item":   &itemController{},
		"battle": &battleController{},
		"order":  &orderController{},
	},
}
