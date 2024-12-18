package backtask

import (
	"html/template"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	"github.com/iufansh/iufans2/taskback"
)

type BackTaskIndexController struct {
	sysmanage.BaseController
}

func (c *BackTaskIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *BackTaskIndexController) Get() {

	list, _ := taskback.GetAllTaskBack()
	c.Data["dataList"] = &list

	c.Data["urlTaskBackIndexGet"] = c.URLFor("BackTaskIndexController.Get")

	if t, err := template.New("tplBackTaskIndex.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
