package jobs

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/task"
	iutils "github.com/iufansh/iutils"
)

func InitSqliteBackup() {
	if dbDriver, _ := config.String("dbdriver"); dbDriver != "sqlite3" {
		return
	}
	backupPath, _ := config.String("sqlite3backuppath")
	if backupPath == "" {
		return
	}
	if !strings.HasSuffix(backupPath, "/") && !strings.HasSuffix(backupPath, "\\") {
		backupPath = backupPath + "/"
	}
	tk1 := task.NewTask("SysSqliteBackup", "15 15 03 * * *", func(ctx context.Context) error {
		logs.Info("SysSqliteBackup start")
		size, err := iutils.CopyFile("./data.db", backupPath+time.Now().Format("20060102150405")+".db")
		if err != nil {
			logs.Error("SysSqliteBackup err:", err)
		} else {
			logs.Info("SysSqliteBackup size =", size)
		}
		files, _ := ioutil.ReadDir(backupPath)
		for i, file := range files {
			if i >= len(files)-7 {
				break
			}
			if err := os.Remove(backupPath + file.Name()); err != nil {
				logs.Error("SysSqliteBackup delete old file =", file.Name(), " err:", err)
			} else {
				logs.Info("SysSqliteBackup delete old file =", file.Name())
			}
		}
		logs.Info("SysSqliteBackup finish")
		return nil
	})
	task.AddTask("SysSqliteBackup", tk1)
	task.StartTask()
}
