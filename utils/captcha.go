package utils

import (
	"strings"

	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/server/web/captcha"
)

var cpt *captcha.Captcha

func InitCaptcha() {
	// use beego cache system store the captcha data
	var domainUri, _ = config.String("domainuri")
	if domainUri != "" && !strings.HasPrefix(domainUri, "/") {
		domainUri = "/" + domainUri
	}
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter(domainUri+"/captcha/", store)
	cpt.ChallengeNums = 4
}

func GetCpt() *captcha.Captcha {
	if cpt == nil {
		InitCaptcha()
	}
	return cpt
}
