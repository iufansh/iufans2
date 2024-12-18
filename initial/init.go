package initial

import (
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/config"
	. "github.com/iufansh/iufans2/utils"
)

func init() {
	InitLog()
	InitSql()
	InitCache()
	InitFilter()
	InitSysTemplateFunc()

	domainUri, _ := config.String("domainuri")
	if domainUri != "" {
		if !strings.HasPrefix(domainUri, "/") {
			domainUri = "/" + domainUri
		}
		web.SetStaticPath(domainUri+"/static", "static")
	}
}
