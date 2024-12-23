package controllers

import (
	"time"

	"github.com/beego/beego/v2/core/logs"
	fu "github.com/iufansh/iufans2/utils"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/orm"
	"github.com/iufansh/iufans2/models"
	license "github.com/iufansh/iuplugins/iu-authorize"
)

type CommonController struct {
	web.Controller
}

// 发送短信验证码，front和api使用 1分钟有效
// 请求：get
// 参数：mobile=手机号&mode=1
// mode=1验证手机号不存在则返回异常；mode=2不验证手机号是否存在；mode=3验证手机号已存在则返回异常
func (c *CommonController) SendSmsCode() {
	logs.Info("\r\n----------request---------",
		"\r\nUri:", c.Ctx.Input.URI(),
		"\r\nMethod:", c.Ctx.Input.Method(),
		"\r\nFrom ip:", c.Ctx.Input.IP(),
		"\r\nUserAgent:", c.Ctx.Input.UserAgent(),
		"\r\nBody:", string(c.Ctx.Input.RequestBody),
		"\r\n--------------------------")
	var code int
	var msg string
	defer func(msg *string, code *int) {
		ret := make(map[string]interface{})
		ret["code"] = code
		ret["msg"] = msg
		ret["data"] = ""
		c.Data["json"] = ret
		c.ServeJSON()
	}(&msg, &code)
	to := c.GetString("mobile")
	if to == "" {
		msg = "手机号不能为空"
		return
	}
	mode, _ := c.GetInt("mode", 1)
	if mode == 1 || mode == 3 { // 验证手机号是否在系统中
		o := orm.NewOrm()
		if exist := o.QueryTable(new(models.Member)).Filter("Username", to).Exist(); mode == 1 && !exist {
			msg = "手机号不存在"
			return
		} else if mode == 3 && exist {
			msg = "手机号已存在"
			return
		}
	}
	sc := models.GetSiteConfigMap(fu.Scsmssignname, fu.Scsmsapi, fu.Scsmsuid, fu.Scsmskey, fu.Scname)
	var companyName string
	if sc[fu.Scsmssignname] != "" {
		companyName = sc[fu.Scsmssignname]
	} else {
		companyName = sc[fu.Scname]
	}
	ms := fu.SmsSender{
		Api:     sc[fu.Scsmsapi],
		Uid:     sc[fu.Scsmsuid],
		Key:     sc[fu.Scsmskey],
		Mobile:  to,
		Company: companyName,
	}
	verifyCode, err := fu.SendSmsVerifyCode(ms)
	var status int
	if err != nil {
		msg = err.Error()
		status = 3
	} else {
		code = 1
		msg = "发送成功"
		status = 2
	}
	// 短信发送记录
	go func(appInfo, receiver, vc, ip string, status int) {
		smsLog := models.SmsLog{
			AppInfo:  appInfo,
			Ip:       ip,
			Receiver: receiver,
			Info:     "验证码：" + vc,
			Status:   status,
		}
		if err := smsLog.InsertLog(); err != nil {
			logs.Error("smsLog.InsertLog err:", err)
		}
	}(c.Ctx.Input.Header("Qx-Api-App"), to, verifyCode, c.Ctx.Input.IP(), status)
}

func (c *CommonController) HealthCheck() {
	if _, err := c.Ctx.ResponseWriter.Write([]byte("1")); err != nil {
		logs.Error("HealthCheck write err:", err)
	}
}

func (c *CommonController) SystemInfo() {
	var code int
	var msg string
	ret := make(map[string]interface{})
	token := c.GetString("t")
	if token != "" {
		t := time.Now().Format("20060102")
		if token == t {
			code = 1
			ret["data"] = license.GetMachineData()
			ret["code"] = code
			ret["msg"] = msg
			c.Data["json"] = ret
			c.ServeJSON()
			return
		}
	}
	c.Abort("404")
}
