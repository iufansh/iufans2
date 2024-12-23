package sysfront

import (
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/iufansh/iufans2/models"
	utils2 "github.com/iufansh/iufans2/utils"
	utils3 "github.com/iufansh/iufans2/utils"
	utils "github.com/iufansh/iutils"
)

type ForgetPwdFrontController struct {
	Base2FrontController
}

func (c *ForgetPwdFrontController) Get() {
	c.TplName = "front/forgetPwd.html"
}

func (c *ForgetPwdFrontController) Post() {
	defer c.RetJSON()
	mobile := strings.TrimSpace(c.GetString("mobile"))
	smsCode := c.GetString("smsCode")
	newPwd := c.GetString("newPwd")
	newPwd2 := c.GetString("newPwd2")
	if mobile == "" || smsCode == "" || newPwd == "" || newPwd2 == "" {
		c.Msg = "信息未填写完整"
		return
	}
	if newPwd != newPwd2 {
		c.Msg = "两次输入的密码不一致"
		return
	}
	if ok := utils2.VerifySmsVerifyCode(mobile, smsCode); !ok {
		c.Msg = "短信验证码错误"
		return
	}
	if !utils3.GetCpt().VerifyReq(c.Ctx.Request) {
		c.Msg = "请刷新验证码，重新输入"
		return
	}
	o := orm.NewOrm()
	member := models.Member{Username: mobile}
	if err := o.Read(&member, "Username"); err != nil {
		c.Msg = "手机号不存在"
		return
	} else {
		salt := utils.GetGuid()
		pa := utils.Md5(utils.Md5(newPwd), utils.Pubsalt, salt)
		member.Password = pa
		member.Salt = salt

		if _, err2 := o.Update(&member, "Password", "Salt", "ModifyDate"); err2 != nil {
			c.Msg = "修改失败，请重试"
			logs.Error("ForgetPwd Change password error", err2)
		} else {
			c.Code = 1
			c.Msg = "修改成功"
			c.Dta = c.URLFor("LoginFrontController.Get")
		}
	}
}
