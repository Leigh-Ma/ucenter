package api

import (
	"ucenter/models"
)

type authorizedController struct {
	apiController
	user *models.User
}

func (c *authorizedController) Prepare() {

}