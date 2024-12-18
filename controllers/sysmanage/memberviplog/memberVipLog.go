package memberviplog

import (
	"html/template"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"
)

type MemberVipLogIndexController struct {
	sysmanage.BaseController
}

func (c *MemberVipLogIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *MemberVipLogIndexController) Get() {
	memberId, _ := c.GetInt64("memberId", -1)

	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := config.Int("pagelimit")
	list, total := new(MemberVipLog).Paginate(page, limit, memberId)
	c.SetPaginator(limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"memberId": memberId}

	c.Data["urlMemberVipLogIndexGet"] = c.URLFor("MemberVipLogIndexController.Get")

	if t, err := template.New("tplIndexMemberVipLog.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}
