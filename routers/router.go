package routers

import (
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/controllers"
	"github.com/iufansh/iufans2/controllers/sysmanage"
	"github.com/iufansh/iufans2/controllers/sysmanage/admin"
	"github.com/iufansh/iufans2/controllers/sysmanage/backtask"
	"github.com/iufansh/iufans2/controllers/sysmanage/gift"
	"github.com/iufansh/iufans2/controllers/sysmanage/index"
	"github.com/iufansh/iufans2/controllers/sysmanage/information"
	"github.com/iufansh/iufans2/controllers/sysmanage/iplist"
	"github.com/iufansh/iufans2/controllers/sysmanage/login"
	"github.com/iufansh/iufans2/controllers/sysmanage/memberlogincount"
	"github.com/iufansh/iufans2/controllers/sysmanage/normalquestion"
	"github.com/iufansh/iufans2/controllers/sysmanage/organization"
	"github.com/iufansh/iufans2/controllers/sysmanage/permission"
	"github.com/iufansh/iufans2/controllers/sysmanage/quicknav"
	"github.com/iufansh/iufans2/controllers/sysmanage/role"
	"github.com/iufansh/iufans2/controllers/sysmanage/siteconfig"
	"github.com/iufansh/iufans2/controllers/sysmanage/smslog"

	"github.com/iufansh/iufans2/controllers/sysapi"
	"github.com/iufansh/iufans2/controllers/sysfront"
	"github.com/iufansh/iufans2/controllers/sysmanage/appversion"
	"github.com/iufansh/iufans2/controllers/sysmanage/member"
	"github.com/iufansh/iufans2/controllers/sysmanage/membersuggest"
	"github.com/iufansh/iufans2/controllers/sysmanage/memberviplog"
	"github.com/iufansh/iufans2/controllers/sysmanage/paymentconfig"
)

func init() {
	// 禁止使用的前缀 i
	web.CtrlGet("/i/:appNo", (*sysfront.AppDownloadFrontController).DownloadRedirect)
	web.CtrlGet("/i/privacy", (*sysfront.PrivacyFrontController).Get)
	web.CtrlGet("/i/privacy/child", (*sysfront.PrivacyFrontController).GetChild)
	web.CtrlGet("i/protocol", (*sysfront.ProtocolFrontController).Get)

	web.CtrlGet("/sendsmscode", (*controllers.CommonController).SendSmsCode) // 通用的短信发送
	web.CtrlGet("/healthcheck", (*controllers.CommonController).HealthCheck)
	web.CtrlGet("/serversysteminfo", (*controllers.CommonController).SystemInfo)

	var adminRouter, _ = config.String("adminrouter")
	web.ErrorController(&controllers.ErrorController{})
	web.CtrlGet(adminRouter+"/sys/main", (*sysmanage.BaseIndexController).Get)
	web.CtrlGet(adminRouter+"/sys/index", (*index.SysIndexController).Get)
	web.CtrlGet(adminRouter+"/sys/getauth", (*index.SysIndexController).GetAuth)
	web.CtrlPost(adminRouter+"/sys/postauth", (*index.SysIndexController).PostAuth)

	web.CtrlPost(adminRouter+"/sys/upload", (*sysmanage.SyscommonController).Upload)
	web.CtrlPost(adminRouter+"/sys/uploadmulti", (*sysmanage.SyscommonController).UploadMulti)

	web.CtrlGet(adminRouter+"/org/index", (*organization.OrganizationIndexController).Get)
	web.CtrlPost(adminRouter+"/org/delone", (*organization.OrganizationIndexController).Delone)
	web.CtrlGet(adminRouter+"/org/add", (*organization.OrganizationAddController).Get)
	web.CtrlPost(adminRouter+"/org/add", (*organization.OrganizationAddController).Post)
	web.CtrlGet(adminRouter+"/org/edit", (*organization.OrganizationEditController).Get)
	web.CtrlPost(adminRouter+"/org/edit", (*organization.OrganizationEditController).Post)

	web.CtrlGet(adminRouter+"/admin/index", (*admin.AdminIndexController).Get)
	web.CtrlGet(adminRouter+"/admin/add", (*admin.AdminAddController).Get)
	web.CtrlPost(adminRouter+"/admin/add", (*admin.AdminAddController).Post)
	web.CtrlGet(adminRouter+"/admin/edit", (*admin.AdminEditController).Get)
	web.CtrlPost(adminRouter+"/admin/edit", (*admin.AdminEditController).Post)
	web.CtrlPost(adminRouter+"/admin/delone", (*admin.AdminIndexController).Delone)
	web.CtrlPost(adminRouter+"/admin/loginverify", (*admin.AdminIndexController).LoginVerify)
	web.CtrlPost(adminRouter+"/admin/locked", (*admin.AdminIndexController).Locked)
	web.CtrlGet(adminRouter+"/changepwd/index", (*admin.ChangePwdController).Get)
	web.CtrlPost(adminRouter+"/changepwd/index", (*admin.ChangePwdController).Post)

	web.CtrlGet(adminRouter+"/role/index", (*role.RoleIndexController).Get)
	web.CtrlPost(adminRouter+"/role/delone", (*role.RoleIndexController).Delone)
	web.CtrlGet(adminRouter+"/role/add", (*role.RoleAddController).Get)
	web.CtrlPost(adminRouter+"/role/add", (*role.RoleAddController).Post)
	web.CtrlGet(adminRouter+"/role/edit", (*role.RoleEditController).Get)
	web.CtrlPost(adminRouter+"/role/edit", (*role.RoleEditController).Post)

	web.CtrlGet(adminRouter+"/permission/index", (*permission.PermissionIndexController).Get)
	web.CtrlPost(adminRouter+"/permission/delone", (*permission.PermissionIndexController).Delone)
	web.CtrlGet(adminRouter+"/permission/add", (*permission.PermissionAddController).Get)
	web.CtrlPost(adminRouter+"/permission/add", (*permission.PermissionAddController).Post)
	web.CtrlGet(adminRouter+"/permission/edit", (*permission.PermissionEditController).Get)
	web.CtrlPost(adminRouter+"/permission/edit", (*permission.PermissionEditController).Post)

	web.CtrlGet(adminRouter+"/login", (*login.LoginController).Get)
	web.CtrlPost(adminRouter+"/login", (*login.LoginController).Post)
	web.CtrlPost(adminRouter+"/loginverify", (*login.LoginController).LoginVerify)
	web.CtrlGet(adminRouter+"/logout", (*login.LoginController).Logout)

	web.CtrlGet(adminRouter+"/site/index", (*siteconfig.SiteConfigIndexController).Get)
	web.CtrlPost(adminRouter+"/site/delone", (*siteconfig.SiteConfigIndexController).Delone)
	web.CtrlGet(adminRouter+"/site/add", (*siteconfig.SiteConfigAddController).Get)
	web.CtrlPost(adminRouter+"/site/add", (*siteconfig.SiteConfigAddController).Post)
	web.CtrlGet(adminRouter+"/site/edit", (*siteconfig.SiteConfigEditController).Get)
	web.CtrlPost(adminRouter+"/site/edit", (*siteconfig.SiteConfigEditController).Post)

	web.CtrlGet(adminRouter+"/information/index", (*information.InformationIndexController).Get)
	web.CtrlPost(adminRouter+"/information/delone", (*information.InformationIndexController).Delone)
	web.CtrlGet(adminRouter+"/information/add", (*information.InformationAddController).Get)
	web.CtrlPost(adminRouter+"/information/add", (*information.InformationAddController).Post)
	web.CtrlGet(adminRouter+"/information/edit", (*information.InformationEditController).Get)
	web.CtrlPost(adminRouter+"/information/edit", (*information.InformationEditController).Post)

	web.CtrlGet(adminRouter+"/normalquestion/index", (*normalquestion.NormalQuestionIndexController).Get)
	web.CtrlPost(adminRouter+"/normalquestion/delone", (*normalquestion.NormalQuestionIndexController).Delone)
	web.CtrlGet(adminRouter+"/normalquestion/add", (*normalquestion.NormalQuestionAddController).Get)
	web.CtrlPost(adminRouter+"/normalquestion/add", (*normalquestion.NormalQuestionAddController).Post)
	web.CtrlGet(adminRouter+"/normalquestion/edit", (*normalquestion.NormalQuestionEditController).Get)
	web.CtrlPost(adminRouter+"/normalquestion/edit", (*normalquestion.NormalQuestionEditController).Post)

	web.CtrlGet(adminRouter+"/qicknav/add", (*quicknav.QuickNavAddController).Get)
	web.CtrlPost(adminRouter+"/qicknav/add", (*quicknav.QuickNavAddController).Post)
	web.CtrlGet(adminRouter+"/qicknav/edit", (*quicknav.QuickNavEditController).Get)
	web.CtrlPost(adminRouter+"/qicknav/edit", (*quicknav.QuickNavEditController).Post)
	web.CtrlGet(adminRouter+"/qicknav/index", (*quicknav.QuickNavIndexController).Get)
	web.CtrlPost(adminRouter+"/qicknav/delone", (*quicknav.QuickNavIndexController).Delone)

	web.CtrlGet(adminRouter+"/iplist/add", (*iplist.IpListAddController).Get)
	web.CtrlPost(adminRouter+"/iplist/add", (*iplist.IpListAddController).Post)
	web.CtrlPost(adminRouter+"/iplist/delone", (*iplist.IpListIndexController).Delone)
	web.CtrlGet(adminRouter+"/iplist/index", (*iplist.IpListIndexController).Get)

	web.CtrlGet(adminRouter+"/paymentconfig/index", (*paymentconfig.PaymentConfigIndexController).Get)
	web.CtrlPost(adminRouter+"/paymentconfig/delone", (*paymentconfig.PaymentConfigIndexController).Delone)
	web.CtrlPost(adminRouter+"/paymentconfig/enabled", (*paymentconfig.PaymentConfigIndexController).Enabled)
	web.CtrlGet(adminRouter+"/paymentconfig/edit", (*paymentconfig.PaymentConfigEditController).Get)
	web.CtrlPost(adminRouter+"/paymentconfig/edit", (*paymentconfig.PaymentConfigEditController).Post)
	web.CtrlGet(adminRouter+"/paymentconfig/add", (*paymentconfig.PaymentConfigAddController).Get)
	web.CtrlPost(adminRouter+"/paymentconfig/add", (*paymentconfig.PaymentConfigAddController).Post)

	web.CtrlGet(adminRouter+"/backtask/index", (*backtask.BackTaskIndexController).Get)
	/* 会员管理 */
	web.CtrlGet(adminRouter+"/member/index", (*member.MemberIndexController).Get)
	web.CtrlPost(adminRouter+"/member/delone", (*member.MemberIndexController).Delone)
	web.CtrlPost(adminRouter+"/member/locked", (*member.MemberIndexController).Locked)
	web.CtrlGet(adminRouter+"/member/edit", (*member.MemberEditController).Get)
	web.CtrlPost(adminRouter+"/member/edit", (*member.MemberEditController).Post)

	web.CtrlPost(adminRouter+"/membersuggest/status", (*membersuggest.MemberSuggestIndexController).Status)
	web.CtrlGet(adminRouter+"/membersuggest/index", (*membersuggest.MemberSuggestIndexController).Get)

	web.CtrlGet(adminRouter+"/memberlogincount/index", (*memberlogincount.MemberLoginCountIndexController).Get)

	web.CtrlGet(adminRouter+"/memberviplog/index", (*memberviplog.MemberVipLogIndexController).Get)

	/* 应用管理 */
	web.CtrlGet(adminRouter+"/appversion/index", (*appversion.AppVersionIndexController).Get)
	web.CtrlPost(adminRouter+"/appversion/delone", (*appversion.AppVersionIndexController).Delone)
	web.CtrlGet(adminRouter+"/appversion/add", (*appversion.AppVersionAddController).Get)
	web.CtrlPost(adminRouter+"/appversion/add", (*appversion.AppVersionAddController).Post)
	web.CtrlGet(adminRouter+"/appversion/edit", (*appversion.AppVersionEditController).Get)
	web.CtrlPost(adminRouter+"/appversion/edit", (*appversion.AppVersionEditController).Post)

	web.CtrlGet(adminRouter+"/gift/index", (*gift.GiftIndexController).Get)
	web.CtrlPost(adminRouter+"/gift/delone", (*gift.GiftIndexController).Delone)
	web.CtrlGet(adminRouter+"/gift/add", (*gift.GiftAddController).Get)
	web.CtrlPost(adminRouter+"/gift/add", (*gift.GiftAddController).Post)

	web.CtrlGet(adminRouter+"/smslog/index", (*smslog.SmsLogIndexController).Get)
	web.CtrlPost(adminRouter+"/smslog/del", (*smslog.SmsLogIndexController).Del)

	// 前端
	var frontRouter, _ = config.String("frontrouter")
	web.CtrlPost(frontRouter+"/upload", (*sysfront.CommonFrontController).Upload)
	web.CtrlGet(frontRouter+"/logout", (*sysfront.LoginFrontController).Logout)
	web.CtrlGet(frontRouter+"/login", (*sysfront.LoginFrontController).Get)
	web.CtrlPost(frontRouter+"/login", (*sysfront.LoginFrontController).Post)
	web.CtrlGet(frontRouter+"/reg", (*sysfront.RegFrontController).Get)
	web.CtrlPost(frontRouter+"/reg", (*sysfront.RegFrontController).Post)
	web.CtrlGet(frontRouter+"/forgetpwd", (*sysfront.ForgetPwdFrontController).Get)
	web.CtrlPost(frontRouter+"/forgetpwd", (*sysfront.ForgetPwdFrontController).Post)
	web.CtrlGet(frontRouter+"/changepwd", (*sysfront.ChangePwdFrontController).Get)
	web.CtrlPost(frontRouter+"/changepwd", (*sysfront.ChangePwdFrontController).Post)

	// api
	var apiRouter, _ = config.String("apirouter")
	web.CtrlPost(apiRouter+"/send/sms", (*sysapi.SendSmsApiController).Post)
	web.CtrlPost(apiRouter+"/loginaliauth", (*sysapi.LoginAliyunAuthApiController).Post)
	web.CtrlPost(apiRouter+"/login", (*sysapi.LoginApiController).Post)
	web.CtrlPost(apiRouter+"/logout", (*sysapi.LoginApiController).Logout)
	web.CtrlPost(apiRouter+"/loginwx", (*sysapi.LoginWxApiController).Post)
	web.CtrlPost(apiRouter+"/loginwxa/userinfo", (*sysapi.LoginWxApiController).PostUserInfo)
	web.CtrlPost(apiRouter+"/loginqq", (*sysapi.LoginQqApiController).Post)
	web.CtrlPost(apiRouter+"/login/apple", (*sysapi.LoginAppleApiController).Post)
	web.CtrlGet(apiRouter+"/loginalipay", (*sysapi.LoginAlipayApiController).Get)
	web.CtrlPost(apiRouter+"/loginalipay", (*sysapi.LoginAlipayApiController).Post)
	web.CtrlPost(apiRouter+"/bindphone", (*sysapi.MemberApiController).BindPhone)
	web.CtrlPost(apiRouter+"/unbindphone", (*sysapi.MemberApiController).UnBindPhone)
	web.CtrlPost(apiRouter+"/cancelaccount", (*sysapi.MemberApiController).CancelAccount)
	web.CtrlPost(apiRouter+"/refreshlogin", (*sysapi.RefreshLoginApiController).Post)
	web.CtrlPost(apiRouter+"/reg", (*sysapi.RegApiController).Post)
	web.CtrlPost(apiRouter+"/tourist/reg", (*sysapi.RegApiController).PostTourist)
	web.CtrlPost(apiRouter+"/forgetpwd", (*sysapi.ForgetPwdApiController).Post)
	web.CtrlPost(apiRouter+"/changepwd", (*sysapi.ChangePwdApiController).Post)
	web.CtrlGet(apiRouter+"/suggest", (*sysapi.MemberSuggestApiController).Get)
	web.CtrlPost(apiRouter+"/suggest", (*sysapi.MemberSuggestApiController).Post)
	web.CtrlGet(apiRouter+"/suggest/unread", (*sysapi.MemberSuggestApiController).GetNewFeedback)
	web.CtrlGet(apiRouter+"/checkupdate", (*sysapi.AppVersionApiController).Get)
	web.CtrlPost(apiRouter+"/checkupdate", (*sysapi.AppVersionApiController).Post)
	web.CtrlPost(apiRouter+"/checkupdate/auto", (*sysapi.AppVersionApiController).PostAuto)
	web.CtrlGet(apiRouter+"/info", (*sysapi.InformationApiController).Get)
	web.CtrlPost(apiRouter+"/info", (*sysapi.InformationApiController).Post)
	web.CtrlGet(apiRouter+"/normalqa", (*sysapi.NormalQuestionApiController).Get)
	web.CtrlGet(apiRouter+"/sysconf", (*sysapi.SysConfigApiController).Get)
	web.CtrlPost(apiRouter+"/member/modifyname", (*sysapi.MemberApiController).ModifyName)
	web.CtrlPost(apiRouter+"/member/uploadavatar", (*sysapi.MemberApiController).UploadAvatar)
	web.CtrlGet(apiRouter+"member/vip", (*sysapi.MemberApiController).GetVip)
}
