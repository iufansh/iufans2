package initial

import (
	"fmt"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
)

func InitLog() {
	filename, err := config.String("logfilename")
	if err != nil {
		filename = "logs/log.log"
	}
	maxdays, err1 := config.Int("logmaxdays")
	level, err2 := config.Int("loglevel")
	if nil != err1 {
		maxdays = 7
	}
	if nil != err2 {
		level = logs.LevelInformational
	}
	logs.SetLogger("file", fmt.Sprintf(`{"filename":"%s","daily":true,"maxdays":%d,"separate":["emergency", "alert", "critical", "error"]}`, filename, maxdays))
	logs.SetLevel(level)
	//logs.SetLogFuncCall(true)
}
