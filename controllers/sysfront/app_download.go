package sysfront

import (
	"net/http"

	"github.com/beego/beego/v2/client/orm"
	"github.com/iufansh/iufans2/models"
)

type AppDownloadFrontController struct {
	Base2FrontController
}

func (c *AppDownloadFrontController) DownloadRedirect() {
	appNo := c.Ctx.Input.Param(":appNo")
	o := orm.NewOrm()
	var appVersion models.AppVersion
	if err := o.QueryTable(new(models.AppVersion)).Filter("AppNo", appNo).Limit(1).OrderBy("-VersionNo").One(&appVersion, "DownloadUrl"); err != nil {
		c.Abort("404")
		return
	}
	c.Redirect(appVersion.DownloadUrl, http.StatusFound)
}
