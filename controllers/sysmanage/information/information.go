package information

import (
	"html/template"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"
)

type InformationIndexController struct {
	sysmanage.BaseController
}

func (c *InformationIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *InformationIndexController) Get() {
	param1 := strings.TrimSpace(c.GetString("param1"))
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := config.Int("pagelimit")
	list, total := new(Information).Paginate(page, limit, c.LoginAdminOrgId, param1)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["urlInformationIndexDelone"] = c.URLFor("InformationIndexController.Delone")
	c.Data["urlInformationAddGet"] = c.URLFor("InformationAddController.Get")
	c.Data["urlInformationEditGet"] = c.URLFor("InformationEditController.Get")

	if t, err := template.New("tplInformationIndex.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *InformationIndexController) Delone() {
	var code int
	var msg string
	url := web.URLFor("InformationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := Information{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		logs.Error("Delete Information eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type InformationAddController struct {
	sysmanage.BaseController
}

func (c *InformationAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *InformationAddController) Get() {
	c.Data["urlInformationIndexGet"] = c.URLFor("InformationIndexController.Get")
	c.Data["urlInformationAddPost"] = c.URLFor("InformationAddController.Post")

	if t, err := template.New("tplAddInformation.tpl").Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *InformationAddController) Post() {
	var code int
	var msg string
	var url = web.URLFor("InformationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := Information{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	receiver := c.GetString("receiver1", "") + ":" + c.GetString("receiver2", "") + ":" + c.GetString("receiver3", "")
	model.Receiver = receiver
	model.Creator = c.LoginAdminId
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Insert(&model); err != nil {
		msg = "添加失败"
		logs.Error("添加失败", err)
	} else {
		code = 1
		msg = "添加成功"
	}
}

type InformationEditController struct {
	sysmanage.BaseController
}

func (c *InformationEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *InformationEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := Information{}
	model.Id = id
	err := o.Read(&model)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(web.URLFor("InformationIndexController.Get"), 302)
	}
	c.Data["receiver"] = strings.Split(model.Receiver, ":")
	c.Data["data"] = &model

	c.Data["urlInformationIndexGet"] = c.URLFor("InformationIndexController.Get")
	c.Data["urlInformationEditPost"] = c.URLFor("InformationEditController.Post")

	if t, err := template.New("tplEditInformation.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplEdit); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *InformationEditController) Post() {
	var code int
	var msg string
	url := web.URLFor("InformationIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := Information{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	receiver := c.GetString("receiver1", "") + ":" + c.GetString("receiver2", "") + ":" + c.GetString("receiver3", "")
	model.Receiver = receiver
	cols := []string{"Title", "Info", "EffectTime", "ExpireTime", "NeedFeedback", "Receiver", "Modifior", "ModifyDate"}
	model.Modifior = c.LoginAdminId
	o := orm.NewOrm()
	if _, err := o.Update(&model, cols...); err != nil {
		msg = "更新失败"
		logs.Error("更新失败", err)
	} else {
		code = 1
		msg = "更新成功"
	}
}
