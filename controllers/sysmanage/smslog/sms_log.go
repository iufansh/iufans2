package smslog

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"
)

type SmsLogIndexController struct {
	sysmanage.BaseController
}

func (c *SmsLogIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *SmsLogIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	timeStart := strings.TrimSpace(c.GetString("timeStart"))
	timeEnd := strings.TrimSpace(c.GetString("timeEnd"))
	status, _ := c.GetInt("status", -1)
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := config.Int("pagelimit")
	list, total := new(SmsLog).Paginate(page, limit, c.LoginAdminOrgId, param1, status, timeStart, timeEnd)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["condArr"] = map[string]interface{}{
		"param1":    param1,
		"timeStart": timeStart,
		"timeEnd":   timeEnd,
		"status":    status,
	}

	c.Data["urlSmsLogIndexGet"] = c.URLFor("SmsLogIndexController.Get")
	c.Data["urlSmsLogIndexDel"] = c.URLFor("SmsLogIndexController.Del")

	if t, err := template.New("tplSmsLogIndex.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *SmsLogIndexController) Del() {
	var code int
	var msg string
	url := web.URLFor("SmsLogIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	param1 := strings.TrimSpace(c.GetString("param1"))
	timeStart := strings.TrimSpace(c.GetString("timeStart"))
	timeEnd := strings.TrimSpace(c.GetString("timeEnd"))
	status, _ := c.GetInt("status", -1)

	num := new(SmsLog).Del(c.LoginAdminOrgId, param1, status, timeStart, timeEnd)
	code = 1
	msg = fmt.Sprintf("成功删除 %d 条", num)
}
