package sysfront

import (
	fm "github.com/iufansh/iufans2/models"
	"github.com/iufansh/iufans2/utils"
)

type ProtocolFrontController struct {
	Base2FrontController
}

func (c *ProtocolFrontController) Get() {
	c.Data["siteName"] = fm.GetSiteConfigValue(utils.Scname)
	c.TplName = "front/protocol.html"
}
