package membersuggest

import (
	"html/template"

	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"

	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
)

type MemberSuggestIndexController struct {
	sysmanage.BaseController
}

func (c *MemberSuggestIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *MemberSuggestIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	status, _ := c.GetInt("status", -1)

	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := config.Int("pagelimit")
	list, total := new(MemberSuggest).Paginate(page, limit, param1, status)
	c.SetPaginator(limit, total)
	// 返回值
	c.Data["dataList"] = &list
	// 查询条件
	c.Data["condArr"] = map[string]interface{}{"param1": param1, "status": status}

	c.Data["urlMemberSuggestIndexGet"] = c.URLFor("MemberSuggestIndexController.Get")
	c.Data["urlMemberSuggestStatus"] = c.URLFor("MemberSuggestIndexController.Status")

	if t, err := template.New("tplIndexMemberSuggest.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *MemberSuggestIndexController) Status() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, err := c.GetInt64("id")
	if err != nil {
		msg = "数据错误"
		logs.Error("Locked MemberSuggest error", err)
		return
	}
	status, _ := c.GetInt("status", 0)
	feedback := c.GetString("feedback")
	o := orm.NewOrm()
	model := MemberSuggest{Id: id}
	if err := o.Read(&model); err != nil {
		logs.Error("Read admin error", err)
		msg = "操作失败，请刷新后重试"
		return
	}
	model.Status = status
	model.Feedback = feedback
	model.Modifior = c.LoginAdminId
	if _, err := o.Update(&model, "Status", "Feedback", "ModifyDate", "Modifior"); err != nil {
		logs.Error("Update MemberSuggest error", err)
		msg = "操作失败，请刷新后重试"
	} else {
		code = 1
		msg = "操作成功"
	}
}
