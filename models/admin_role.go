package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type AdminRole struct {
	Id      int64
	AdminId int64
	RoleId  int64
}

func init() {
	orm.RegisterModelWithPrefix(SysDbPrefix, new(AdminRole))
}

func (model *AdminRole) TableUnique() [][]string {
	return [][]string{
		{"AdminId", "RoleId"},
	}
}
