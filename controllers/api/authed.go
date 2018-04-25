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

func (c *authorizedController) currentPlayer() *models.Player {
	return models.GetPlayer(c.user.GetId())
}