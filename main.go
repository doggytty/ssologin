package main

import (
	_ "github.com/doggytty/ssologin/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/doggytty/ssologin/models"
	"github.com/doggytty/ssologin/filters"
)

func main() {
	// beego默认logger
	logs.Async()
	// log异步输出
	logs.EnableFuncCallDepth(true)
	if beego.BConfig.RunMode == "dev" {
		logs.SetLogger(logs.AdapterConsole)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"login.log","daily":true,"maxdays":10}`)
	}

	// init database
	models.SyncDataBase()

	// 修改模板的位置
	beego.BConfig.WebConfig.ViewsPath = "templates"
	// 关闭自动渲染
	//beego.BConfig.WebConfig.AutoRender = false
	// {{ 和 }} 作为左右标签
	//beego.BConfig.WebConfig.TemplateLeft = "{{"
	//beego.BConfig.WebConfig.TemplateRight = "}}"
	// 自定义模板后缀名称
	beego.AddTemplateExt(".html")
	// 默认的静态文件处理路径
	beego.SetStaticPath("/static", "static")

	// 启用session
	beego.BConfig.WebConfig.Session.SessionOn = true
	//beego.BConfig.WebConfig.Session.SessionProvider = "memory"
	//beego.BConfig.WebConfig.Session.SessionProvider = "redis"
	//beego.BConfig.WebConfig.Session.SessionProviderConfig = "127.0.0.1:6379"

	// filter
	// BeforeStatic 静态地址之前
	// BeforeRouter 寻找路由之前
	// BeforeExec 找到路由之后，开始执行相应的 Controller 之前
	// AfterExec 执行完 Controller 逻辑之后执行的过滤器
	// FinishRouter 执行完逻辑之后执行的过滤器
	// 使用session必须在beforeStatic之后
	// 1、黑/白名单、DDOS防止
	beego.InsertFilter("/*", beego.BeforeStatic, filters.FilterBlackDDOS)
	// 2、session filter
	beego.InsertFilter("/*", beego.BeforeRouter, filters.FilterLogin)
	// 3、admin filter
	beego.InsertFilter("/admin/*", beego.BeforeRouter, filters.FilterAdministrator)
	// 启动beego
	beego.Run()
}

