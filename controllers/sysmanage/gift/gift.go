package gift

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"
)

type GiftIndexController struct {
	sysmanage.BaseController
}

func (c *GiftIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *GiftIndexController) Get() {

	param1 := c.GetString("param1")
	status, _ := c.GetInt("status", -1)
	page, err := c.GetInt("p")
	if err != nil {
		page = 1
	}
	limit, _ := config.Int("pagelimit")
	list, total := new(Gift).Paginate(page, limit, c.LoginAdminOrgId, param1, status)
	c.SetPaginator(limit, total)
	c.Data["dataList"] = &list

	c.Data["condArr"] = map[string]interface{}{"param1": param1, "status": status}

	c.Data["urlGiftIndexDelone"] = c.URLFor("GiftIndexController.Delone")
	c.Data["urlGiftAddGet"] = c.URLFor("GiftAddController.Get")
	c.Data["urlGiftEditGet"] = c.URLFor("GiftEditController.Get")

	if t, err := template.New("tplGiftIndex.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *GiftIndexController) Delone() {
	var code int
	var msg string
	url := web.URLFor("GiftIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	model := Gift{}
	model.Id = id
	err := o.Read(&model)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	}
	_, err1 := o.Delete(&model, "Id")
	if err1 != nil {
		logs.Error("Delete Gift eror", err1)
		msg = "删除失败"
	} else {
		code = 1
		msg = "删除成功"
	}
}

type GiftAddController struct {
	sysmanage.BaseController
}

func (c *GiftAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *GiftAddController) Get() {
	c.Data["urlGiftIndexGet"] = c.URLFor("GiftIndexController.Get")
	c.Data["urlGiftAddPost"] = c.URLFor("GiftAddController.Post")
	if t, err := template.New("tplAddGift.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"urlfor": web.URLFor,
	}).Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *GiftAddController) Post() {
	var code int
	var msg string
	var url = web.URLFor("GiftIndexController.Get")
	defer sysmanage.Retjson(c.Ctx, &msg, &code, &url)
	model := Gift{}
	if err := c.ParseForm(&model); err != nil {
		msg = "参数异常"
		return
	}
	codes := strings.Split(model.Code, "\r\n")
	models := make([]Gift, 0)
	for _, v := range codes {
		if strings.TrimSpace(v) == "" {
			continue
		}
		models = append(models, Gift{
			OrgId:      c.LoginAdminOrgId,
			CreateDate: time.Now(),
			Creator:    c.LoginAdminId,
			AppNo:      model.AppNo,
			Price:      model.Price,
			Code:       strings.TrimSpace(v),
			Status:     0,
		})
	}
	o := orm.NewOrm()
	if num, err := o.InsertMulti(len(models), models); err != nil {
		msg = "添加失败"
		return
	} else {
		msg = fmt.Sprintf("成功添加 %d 条", num)
		code = 1
	}
}
