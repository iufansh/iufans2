package role

var tplEdit = `
<!DOCTYPE html>
<html>
<head>
{{.HtmlHead}}
</head>
<body>
<div class="layui-fluid">
    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs12 layui-col-sm12 layui-col-md12">
			<!--tab标签-->
			<div class="layui-tab layui-tab-brief">
				<ul class="layui-tab-title">
					<li class=""><a href='{{.urlRoleIndexGet}}'>角色列表</a></li>
					<li class="layui-this">编辑角色</li>
				</ul>
				<div class="layui-tab-content">
					<div class="layui-tab-item layui-show">
						<form class="layui-form form-container" action='{{.urlRoleEditPost}}' method="post">
							{{.xsrfdata}}
							<input type="hidden" name="Id" value="{{.data.Id}}" >
							<div class="layui-form-item">
								<label class="layui-form-label">角色名称</label>
								<div class="layui-input-block">
									<input type="text" name="Name" value="{{.data.Name}}" required lay-verify="required" placeholder="请输入角色名称" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">描述</label>
								<div class="layui-input-block">
									<input type="text" name="Description" value="{{.data.Description}}" placeholder="请输入描述" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">默认首页</label>
								<div class="layui-input-block">
									<input type="text" name="HomeUrl" value="{{.data.HomeUrl}}" placeholder="登陆后显示的首页，格式可 /home或HomeController.Get" class="layui-input">
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">状态</label>
								<div class="layui-input-block">
									<input type="radio" name="Enabled" value="1" title="启用" {{if eq .data.Enabled 1}}checked="checked"{{end}}>
									<input type="radio" name="Enabled" value="0" title="禁用" {{if eq .data.Enabled 0}}checked="checked"{{end}}>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">开放给组织</label>
								<div class="layui-input-block">
									<input type="radio" name="IsOrg" value="1" title="启用" {{if eq .data.IsOrg 1}}checked="checked"{{end}}>
									<input type="radio" name="IsOrg" value="0" title="禁用" {{if eq .data.IsOrg 0}}checked="checked"{{end}}>
								</div>
							</div>
							<div class="layui-form-item">
								<label class="layui-form-label">权限</label>
								{{range $index, $vo := .permissionList}}
									{{if eq $vo.Pid 0}}
									<div class="layui-input-block">
										<input type="checkbox" name="permissions" title="{{$vo.Name}}" value="{{$vo.Id}}" {{if index $.rolePermissionMap $vo.Id}} checked {{end}} lay-skin="primary">
									</div>
										{{range $index1, $vo1 := $.permissionList}}
											{{if eq $vo.Id $vo1.Pid}}
											<div class="layui-input-block">
												&nbsp;&nbsp;|__
												<input type="checkbox" name="permissions" title="{{$vo1.Name}}" value="{{$vo1.Id}}" {{if index $.rolePermissionMap $vo1.Id}} checked {{end}} lay-skin="primary">
											</div>
												{{range $index2, $vo2 := $.permissionList}}
													{{if eq $vo1.Id $vo2.Pid}}
													<div class="layui-input-block">
														&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|__
														<input type="checkbox" name="permissions" title="{{$vo2.Name}}" value="{{$vo2.Id}}" {{if index $.rolePermissionMap $vo2.Id}} checked {{end}} lay-skin="primary">
													</div>
														{{range $index3, $vo3 := $.permissionList}}
															{{if eq $vo2.Id $vo3.Pid}}
															<div class="layui-input-block">
																&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;|__
																<input type="checkbox" name="permissions" title="{{$vo3.Name}}" value="{{$vo3.Id}}" {{if index $.rolePermissionMap $vo3.Id}} checked {{end}} lay-skin="primary">
															</div>
															{{end}}
														{{end}}
													{{end}}
												{{end}}
											{{end}}
										{{end}}
									{{end}}
								{{else}}
								<div class="layui-input-block">
									<label class="layui-form-label">未配置菜单</label>
								</div>
								{{end}}
							</div>
							<div class="layui-form-item">
								<div class="layui-input-block">
									<button class="layui-btn" lay-submit lay-filter="*">保存</button>
								</div>
							</div>
						</form>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
{{.Scripts}}
</body>
</html>
`
