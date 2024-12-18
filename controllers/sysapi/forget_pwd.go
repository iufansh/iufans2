package sysapi

import (
	"encoding/json"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/iufansh/iufans2/models"
	"github.com/iufansh/iufans2/utils"
	"github.com/iufansh/iutils"
)

type forgetPwdParam struct {
	Mobile   string `json:"mobile"`   // 必填
	AuthCode string `json:"authCode"` // 必填
}

type ForgetPwdApiController struct {
	Base2ApiController
}

/*
api验证短信
param:
body: {"mobile":"13112345678","authCode":3256}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","phone":"13111111111","nickname":"昵称","autoLogin":true}}
desc: 验证短信验证码是否正确
*/
func (c *ForgetPwdApiController) Post() {
	defer c.RetJSON()
	var p forgetPwdParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	if p.Mobile == "" {
		c.Msg = "手机号不能为空"
		return
	}
	if p.AuthCode == "" {
		c.Msg = "验证码不能为空"
		return
	}
	if ok := utils.VerifySmsVerifyCode(p.Mobile, p.AuthCode); !ok {
		c.Msg = "短信验证码错误"
		return
	}
	o := orm.NewOrm()
	member := models.Member{Username: p.Mobile}
	if err := o.Read(&member, "Username"); err != nil {
		if err == orm.ErrNoRows {
			c.Msg = "手机号不存在"
			return
		} else {
			c.Msg = "验证失败，请重试"
			return
		}
	}
	// 自动登录
	member.LoginIp = c.Ctx.Input.IP()
	// 以下3个是用于统计登录次数
	member.AppNo = c.AppNo
	member.AppChannel = c.AppChannel
	member.AppVersion = c.AppVersionCode
	_, _, token := UpdateMemberLoginStatus(member)

	c.Code = utils.CODE_OK
	c.Msg = "验证成功"
	var vipEffect int
	if member.Vip > 0 && !member.VipExpire.IsZero() && member.VipExpire.After(time.Now().AddDate(0, 0, -1)) {
		vipEffect = 1
	}
	c.Dta = map[string]interface{}{
		"id":         member.Id,
		"token":      token,
		"phone":      member.GetFmtMobile(),
		"nickname":   member.Name,
		"autoLogin":  true,
		"avatar":     member.GetFullAvatar(c.Ctx.Input.Site()),
		"inviteCode": utils.GenInviteCode(member.Id),
		"vipEffect":  vipEffect,
		"vip":        member.Vip,
		"vipExpire":  iutils.FormatDate(member.VipExpire),
	}
}
