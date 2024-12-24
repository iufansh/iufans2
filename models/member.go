package models

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
)

type Member struct {
	Id                int64     `orm:"auto"` // 自增主键
	CreateDate        time.Time `orm:"auto_now_add;type(datetime)" description:"创建时间"`
	ModifyDate        time.Time `orm:"auto_now;type(datetime)" description:"更新时间"`
	Creator           int64     `description:"创建人ID"`
	Modifior          int64     `description:"更新人ID"`
	Version           int       `description:"版本"`
	OrgId             int64     `description:"组织ID"`  // 组织ID
	RefId             int64     `description:"推荐人ID"` // 推荐人ID
	Levels            string    `description:"层级关系"`  // 层级关系
	LevelsDeep        int       `description:"层级深度"`  // 层级深度
	AppNo             string    `description:"App编号"`
	AppChannel        string    `description:"App渠道"`
	AppVersion        int       `description:"App版本"`
	Username          string    `orm:"unique;size(127)" description:"用户名"`
	ThirdAuthId       string    `description:"三方登录ID"` // 三方登录的ID, 比如微信的unionid，华为的AuthHuaweiId
	RegType           int       `description:"注册类型"`   // 注册类型 0-系统创建；1-手机号；2-微信；3-支付宝；4-QQ；5-本机号码一键登录；6-Apple登录；7-游客模式
	Name              string    `description:"昵称"`
	Mobile            string    `description:"手机号"`
	Password          string    `description:"密码"`
	Salt              string    `description:"密码加盐"`
	Vip               int       `description:"是否VIP"`
	VipTime           time.Time `orm:"null" description:"最近VIP获得时间"` // 最近VIP获得时间
	VipExpire         time.Time `orm:"null" description:"VIP过期时间"`   // VIP过期时间
	Avatar            string    `description:"头像"`
	Locked            int8      `description:"是否锁定"`
	LockedDate        time.Time `orm:"null" description:"锁定时间"`
	LoginDate         time.Time `orm:"null" description:"登录时间"`
	LoginFailureCount int       `description:"登录失败数"`
	LoginIp           string    `description:"登录IP"`
	Enabled           int8      `description:"是否可用"`
	Token             string    `description:"Token"`
	TokenExpTime      time.Time `orm:"null" description:"Token过期时间"`
	Cancelled         int8      `description:"是否注销"` // 是否注销 0-正常使用；1-已注销
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(Member))
}

/**
 * 获取格式化的手机号，格式：131*****234
 */
func (model *Member) GetFmtMobile() string {
	if len(model.Mobile) != 11 {
		return ""
	}
	return web.Substr(model.Mobile, 0, 3) + "*****" + web.Substr(model.Username, 8, 3)
}

/*
 * 获取完整的头像地址
 */
func (model *Member) GetFullAvatar(domain string) string {
	if model.Avatar == "" {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(model.Avatar), "http") {
		return model.Avatar
	}
	if strings.HasPrefix(model.Avatar, "/") {
		return domain + model.Avatar
	}
	return domain + "/" + model.Avatar
}

func (model *Member) Paginate(page int, limit int, orderBy int, id int64, param1 string, regType int) (list []Member, total int64) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	o := orm.NewOrm()
	qs := o.QueryTable(new(Member))
	cond := orm.NewCondition()
	if param1 != "" {
		cond = cond.AndCond(cond.And("Name__contains", param1).Or("Username__contains", param1).Or("Mobile__contains", param1))
	}
	if regType != 0 {
		cond = cond.And("RegType", regType)
	}
	if id != -1 {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("RefId", id).Or("Id", id)
		cond = cond.AndCond(cond1)
	}
	qs = qs.SetCond(cond)
	qs = qs.Limit(limit)
	qs = qs.Offset(offset)
	switch orderBy {
	case 1:
		qs = qs.OrderBy("Id")
		break
	case 2:
		qs = qs.OrderBy("-LoginDate")
		break
	case 3:
		qs = qs.OrderBy("LoginDate")
		break
	default:
		qs = qs.OrderBy("-Id")
		break
	}
	qs.All(&list)
	total, _ = qs.Count()
	return
}
