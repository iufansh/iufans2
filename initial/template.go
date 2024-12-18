package initial

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
	"github.com/iufansh/iufans2/utils"

	"strings"

	utils2 "github.com/iufansh/iutils"
)

func InitSysTemplateFunc() {
	web.AddFuncMap("getSiteConfigCodeMap", utils.GetSiteConfigCodeMap)
	web.AddFuncMap("numberAdd", numberAdd)
	web.AddFuncMap("numberMinus", numberMinus)
	web.AddFuncMap("formatAmount", formatAmount)
	web.AddFuncMap("list2Map", utils2.List2Map)
	web.AddFuncMap("date", formatDate)
	web.AddFuncMap("urlfor", web.URLFor)
	utils.InitCaptcha()
}

// 自定义 Date 函数
func formatDate(t time.Time, format string) string {
    // 替换 Beego 风格的格式为标准库格式
    format = strings.ReplaceAll(format, "Y", "2006")
    format = strings.ReplaceAll(format, "m", "01")
    format = strings.ReplaceAll(format, "d", "02")
    format = strings.ReplaceAll(format, "H", "15")
    format = strings.ReplaceAll(format, "i", "04")
    format = strings.ReplaceAll(format, "s", "05")

    return t.Format(format)
}

// 整数相减
func numberMinus(a, b interface{}) int64 {
	return numberAdd(a, 0) - numberAdd(0, b)
}

// 整数相加
func numberAdd(a, b interface{}) int64 {
	var aint64 int64
	var bint64 int64
	switch a.(type) {
	case int64:
		aint64 = a.(int64)
		break
	case int32:
		aint64 = int64(a.(int32))
		break
	case int:
		aint64 = int64(a.(int))
		break
	case string:
		if i, err := strconv.ParseInt(a.(string), 10, 64); err == nil {
			aint64 = i
			break
		}
	case int16:
		aint64 = int64(a.(int16))
		break
	case int8:
		aint64 = int64(a.(int8))
		break
	}
	switch b.(type) {
	case int64:
		bint64 = b.(int64)
		break
	case int32:
		bint64 = int64(b.(int32))
		break
	case int:
		bint64 = int64(b.(int))
		break
	case string:
		if i, err := strconv.ParseInt(b.(string), 10, 64); err == nil {
			bint64 = i
			break
		}
	case int16:
		bint64 = int64(b.(int16))
		break
	case int8:
		bint64 = int64(b.(int8))
		break
	}
	return aint64 + bint64
}

func formatAmount(a interface{}) string {
	var f float64
	switch a.(type) {
	case int64:
		f = float64(a.(int64)) / 100
	case int32:
		f = float64(int64(a.(int32))) / 100
	case int:
		f = float64(a.(int)) / 100
	case string:
		if i, err := strconv.ParseFloat(a.(string), 64); err == nil {
			f = i / 100
		}
	case float32:
		f = float64(a.(float32)) / 100
	case float64:
		f = a.(float64) / 100
	default:
		f = -1.00
	}
	s := strconv.FormatFloat(f, 'f', 2, 64)
	return strings.TrimSuffix(s, ".00")
}
