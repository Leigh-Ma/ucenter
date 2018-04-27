package user

import "ucenter/controllers"

var UserRouter = &controllers.RouterGroup{
	Namespace: "/api/",
	Routers: map[string]controllers.IExport{
		"login":    &loginController{},
		"token":    &tokenController{},
		"user":     &userController{},
		"register": &registerController{},
		"order":    &orderCbController{},
	},
}
