package role

import (
	"context"
	"html/template"
	"strconv"

	"github.com/iufansh/iufans2/controllers/sysmanage"
	. "github.com/iufansh/iufans2/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	"github.com/beego/beego/v2/server/web"
)

func validate(role *Role) (hasError bool, errMsg string) {
	valid := validation.Validation{}
	valid.Required(role.Name, "errmsg").Message("角色名必输")
	valid.MaxSize(role.Name, 50, "errmsg").Message("角色名最长50位")
	if valid.HasErrors() {
		for _, err := range valid.Errors {
			return true, err.Message
		}
	}
	return false, ""
}

type RoleIndexController struct {
	sysmanage.BaseController
}

func (c *RoleIndexController) NestPrepare() {
	c.EnableRender = false
}

func (c *RoleIndexController) Get() {
	var roleList []Role
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	qs.All(&roleList)
	// 返回值
	c.Data["dataList"] = &roleList

	c.Data["urlRoleIndexDelone"] = c.URLFor("RoleIndexController.Delone")
	c.Data["urlRoleAddGet"] = c.URLFor("RoleAddController.Get")
	c.Data["urlRoleEditGet"] = c.URLFor("RoleEditController.Get")

	if t, err := template.New("tplPermissionIndex.tpl").Funcs(map[string]interface{}{ // 这个模式加载的模板，必须在这里注册模板函数，无法使用内置的模板函数
		"date": web.Date,
	}).Parse(tplIndex); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *RoleIndexController) Delone() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	id, _ := c.GetInt64("id")
	role := Role{Id: id}
	om := orm.NewOrm()
	err := om.Read(&role)
	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		code = 1
		msg = "删除成功"
		return
	} else if role.IsSystem == 1 {
		msg = "无法删除"
		return
	}
	// 先删除角色权限关联
	om.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		if _, err := txOrm.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
			logs.Error("Delete role error 1", err)
			msg = "删除失败"
			return err
		}

		if _, err := txOrm.Delete(&Role{Id: id}); err != nil {
			logs.Error("Delete role error 2", err)
			msg = "删除失败"
			return err
		}
		code = 1
		msg = "删除成功"
		return nil
	})

}

type RoleAddController struct {
	sysmanage.BaseController
}

func (c *RoleAddController) NestPrepare() {
	c.EnableRender = false
}

func (c *RoleAddController) Get() {
	c.Data["permissionList"] = GetPermissionList()

	c.Data["urlRoleIndexGet"] = c.URLFor("RoleIndexController.Get")
	c.Data["urlRoleAddPost"] = c.URLFor("RoleAddController.Post")

	if t, err := template.New("tplAddRole.tpl").Parse(tplAdd); err != nil {
		logs.Error("template Parse err", err)
	} else {
		t.Execute(c.Ctx.ResponseWriter, c.Data)
	}
}

func (c *RoleAddController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	role := Role{}
	if err := c.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	} else if hasError, errMsg := validate(&role); hasError {
		msg = errMsg
		return
	}
	role.Creator = c.LoginAdminId
	role.Modifior = c.LoginAdminId
	om := orm.NewOrm()
	om.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		if _, err := txOrm.Insert(&role); err != nil {
			msg = "添加失败"
			logs.Error("Insert role error 1", err)
			return err
		} else {
			permissions := c.GetStrings("permissions")
			rolePermissions := make([]RolePermission, 0)
			for _, v := range permissions {
				permissionId, _ := strconv.ParseInt(v, 10, 64)
				ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
				rolePermissions = append(rolePermissions, ar)
			}
			if _, err := txOrm.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
				msg = "添加失败"
				logs.Error("Insert role error 2", err)
				return err
			}
			code = 1
			msg = "添加成功"
			return nil
		}
	})

}

type RoleEditController struct {
	sysmanage.BaseController
}

func (c *RoleEditController) NestPrepare() {
	c.EnableRender = false
}

func (c *RoleEditController) Get() {
	id, _ := c.GetInt64("id")
	o := orm.NewOrm()
	role := Role{Id: id}

	err := o.Read(&role)

	if err == orm.ErrNoRows || err == orm.ErrMissPK {
		c.Redirect(web.URLFor("RoleIndexController.get"), 302)
	} else {
		// 当前角色包含的权限
		var rpList orm.ParamsList
		o.QueryTable(new(RolePermission)).Filter("RoleId", id).ValuesFlat(&rpList, "PermissionId")
		rpMap := make(map[int64]bool)
		for _, v := range rpList {
			rpId, ok := v.(int64)
			if ok {
				rpMap[rpId] = true
			}
		}
		c.Data["data"] = &role
		c.Data["rolePermissionMap"] = rpMap
		c.Data["permissionList"] = GetPermissionList()

		c.Data["urlRoleIndexGet"] = c.URLFor("RoleIndexController.Get")
		c.Data["urlRoleEditPost"] = c.URLFor("RoleEditController.Post")

		if t, err := template.New("tplEditRole.tpl").Parse(tplEdit); err != nil {
			logs.Error("template Parse err", err)
		} else {
			t.Execute(c.Ctx.ResponseWriter, c.Data)
		}
	}
}

func (c *RoleEditController) Post() {
	var code int
	var msg string
	defer sysmanage.Retjson(c.Ctx, &msg, &code)
	role := Role{}
	if err := c.ParseForm(&role); err != nil {
		msg = "参数异常"
		return
	}
	cols := []string{"Name", "Description", "HomeUrl", "Enabled", "ModifyDate"}
	role.Modifior = c.LoginAdminId
	om := orm.NewOrm()
	om.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		if _, err := txOrm.Update(&role, cols...); err != nil {
			msg = "更新失败"
			logs.Error("Update role error 1", err)
			return err
		} else {
			// 删除旧权限
			if _, err := txOrm.QueryTable(new(RolePermission)).Filter("RoleId", role.Id).Delete(); err != nil {
				msg = "更新失败"
				logs.Error("Update role error 2", err)
				return err
			}
			// 重新插入新权限
			permissions := c.GetStrings("permissions")
			rolePermissions := make([]RolePermission, 0)
			for _, v := range permissions {
				permissionId, _ := strconv.ParseInt(v, 10, 64)
				ar := RolePermission{RoleId: role.Id, PermissionId: permissionId}
				rolePermissions = append(rolePermissions, ar)
			}
			if _, err := txOrm.InsertMulti(len(rolePermissions), rolePermissions); err != nil {
				msg = "更新失败"
				logs.Error("Update role error 3", err)
				return err
			}
			code = 1
			msg = "更新成功"
			return nil
		}
	})

}
