package sysapi

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	fm "github.com/iufansh/iufans2/models"
	. "github.com/iufansh/iutils"

	"encoding/json"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/validation"
	. "github.com/iufansh/iufans2/models"
	"github.com/iufansh/iufans2/utils"
)

type RegApiController struct {
	Base2ApiController
}

type regParam struct {
	InviteCode string `json:"inviteCode"`
	AuthCode   string `json:"authCode"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

/*
api注册
param:
body:{"inviteCode":3506,"authCode":2356,"username":"aaaaa","password":"32md5小写"}
return:{"code":1,"msg":"成功","data":{"id":1,"token":"11111111111111111111","phone":"13111111111","nickname":"昵称","autoLogin":true}}
desc:注册成功，记录登录状态
*/
func (c *RegApiController) Post() {
	defer c.RetJSON()
	var p regParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	p.Password = strings.ToLower(p.Password)
	model := Member{
		Username: p.Username,
		Password: p.Password,
	}

	scMap := fm.GetSiteConfigMap(utils.Scfrontregsmsverify)
	if scMap[utils.Scfrontregsmsverify] == "1" {
		if p.AuthCode == "" {
			c.Msg = "短信验证码必填"
			return
		}
		if ok := utils.VerifySmsVerifyCode(model.Username, p.AuthCode); !ok {
			c.Msg = "短信验证码错误"
			return
		}
	}
	valid := validation.Validation{}
	valid.Required(model.Username, "errmsg").Message("手机号必填")
	valid.MaxSize(model.Username, 11, "errmsg").Message("手机号最长11位")
	valid.Required(model.Password, "errmsg").Message("密码必填")
	valid.Length(model.Password, 32, "errmsg").Message("密码格式错误")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			c.Msg = err.Message
			return
		}
	}
	// 验证用户名是否存在
	o := orm.NewOrm()
	if isExist := o.QueryTable(new(Member)).Filter("Username", model.Username).Exist(); isExist {
		c.Msg = "账号已存在"
		return
	}
	// 查询层级
	var refId = utils.ReverseInviteCode(p.InviteCode)

	model, err := CreateMemberReg(1, c.AppNo, c.AppChannel, c.AppVersionCode, refId, model.Username, model.Password, model.Name, "", "")
	if err != nil {
		c.Msg = "注册失败，请重试"
		return
	}
	// 自动登录
	model.LoginIp = c.Ctx.Input.IP()
	// 以下两个是用于统计登录次数
	model.AppNo = c.AppNo
	model.AppChannel = c.AppChannel
	model.AppVersion = c.AppVersionCode
	_, _, token := UpdateMemberLoginStatus(model)

	c.Msg = "注册成功"
	c.Code = utils.CODE_OK
	var vipEffect int
	if model.Vip > 0 && !model.VipExpire.IsZero() && model.VipExpire.After(time.Now().AddDate(0, 0, -1)) {
		vipEffect = 1
	}
	c.Dta = map[string]interface{}{
		"id":         model.Id,
		"token":      token,
		"phone":      model.GetFmtMobile(),
		"nickname":   model.Name,
		"autoLogin":  true,
		"avatar":     model.GetFullAvatar(c.Ctx.Input.Site()),
		"inviteCode": utils.GenInviteCode(model.Id),
		"vipEffect":  vipEffect,
		"vip":        model.Vip,
		"vipExpire":  FormatDate(model.VipExpire),
	}
	// go GenerateRandAvatar(memberId)
	return
}

type touristParam struct {
	Uuid string `json:"uuid"`
}

func (c *RegApiController) PostTourist() {
	defer c.RetJSON()
	var p touristParam
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &p); err != nil {
		c.Code = utils.CODE_ERROR
		c.Msg = "参数格式错误"
		return
	}
	o := orm.NewOrm()
	var model Member
	var isNew bool
	if p.Uuid == "" {
		p.Uuid = "TOU_" + Md5("tourist", strconv.FormatInt(time.Now().Unix(), 10), strconv.FormatInt(int64(RandNum(1, 1000000)), 10))
		isNew = true
	} else {
		if err := o.QueryTable(new(Member)).Filter("Username", p.Uuid).One(&model); err != nil {
			if err == orm.ErrNoRows {
				isNew = true
			} else {
				c.Msg = "访问失败，请重试"
				return
			}
		}
	}
	if isNew {
		var err error
		model, err = CreateMemberReg(7, c.AppNo, c.AppChannel, c.AppVersionCode, 0, p.Uuid, p.Uuid, "", "", "")
		if err != nil {
			c.Msg = "访问失败，请重试W2"
			return
		}
	}

	// 自动登录
	model.LoginIp = c.Ctx.Input.IP()
	// 以下两个是用于统计登录次数
	model.AppNo = c.AppNo
	model.AppChannel = c.AppChannel
	model.AppVersion = c.AppVersionCode
	_, _, token := UpdateMemberLoginStatus(model)

	c.Msg = "访问成功"
	c.Code = utils.CODE_OK
	var vipEffect int
	if model.Vip > 0 && !model.VipExpire.IsZero() && model.VipExpire.After(time.Now().AddDate(0, 0, -1)) {
		vipEffect = 1
	}
	c.Dta = map[string]interface{}{
		"uuid":       model.Username, // 仅游客模式有
		"id":         model.Id,
		"token":      token,
		"phone":      model.GetFmtMobile(),
		"nickname":   model.Name,
		"autoLogin":  true,
		"avatar":     model.GetFullAvatar(c.Ctx.Input.Site()),
		"inviteCode": utils.GenInviteCode(model.Id),
		"vipEffect":  vipEffect,
		"vip":        model.Vip,
		"vipExpire":  FormatDate(model.VipExpire),
	}
	return
}

func CreateMemberReg(regType int, appNo, appChannel string, appVersion int, refId int64, username string, password string, name string, thirdAuthId string, avatar string) (model Member, err error) {
	model.RegType = regType
	model.AppNo = appNo
	model.AppChannel = appChannel
	model.AppVersion = appVersion
	model.RefId = refId
	model.Username = username
	model.ThirdAuthId = thirdAuthId
	//if avatar == "" {
	//	model.Avatar = "/static/front/images/avatar/default.png"
	//} else {
	model.Avatar = avatar
	//}
	if name == "" {
		if len(model.Username) == 11 && strings.HasPrefix(model.Username, "1") {
			model.Name = SubString(model.Username, 0, 3) + "*****" + SubString(model.Username, 8, 3)
		} else {
			if model.RegType == 7 {
				model.Name = fmt.Sprintf("游客%d", RandNum(1, 1000000))
			} else {
				model.Name = fmt.Sprintf("会员%d", RandNum(1, 1000000))
			}
		}
	} else {
		model.Name = name
	}
	if len(model.Username) == 11 && strings.HasPrefix(model.Username, "1") {
		model.Mobile = model.Username
	}
	salt := GetGuid()
	model.Password = Md5(password, Pubsalt, salt)
	model.Salt = salt
	model.Creator = 0
	model.Modifior = 0
	model.Enabled = 1
	model.Locked = 0
	model.LoginFailureCount = 0

	var memberId int64
	o := orm.NewOrm()
	if memberId, err = o.Insert(&model); err != nil {
		logs.Error("memberRegErr Member error", err)
		return model, err
	}
	model.Id = memberId
	//缓存头像到本地
	// go GetMemberAvatar(memberId, model.Avatar)
	return model, nil
}

// 获取头像
func GetMemberAvatar(id int64, avatar string) {
	//logs.Info(avatar)
	var avatarFile string
	// 网络图片，下载到本地缓存
	if avatar != "" && strings.HasPrefix(strings.ToLower(avatar), "http") {
		avatarFile = fmt.Sprintf("upload/avatar/%d.jpg", id)
		if err := httplib.Get(avatar).ToFile(avatarFile); err != nil {
			logs.Error("GetMemberAvatar err:", err)
			return
		}
		avatarFile = "/" + avatarFile
	} else {
		if strings.HasPrefix(avatar, "/") {
			avatarFile = avatar
		} else {
			avatarFile = "/" + avatar
		}
	}
	o := orm.NewOrm()
	o.Update(&Member{Id: id, Avatar: avatarFile}, "Avatar")
}
