package routers

import (
	"ucenter/controllers/api"
	"ucenter/controllers/user"
)

func init() {
	api.ApiRouter.RegisterRouter()
	user.UserRouter.RegisterRouter()
}
