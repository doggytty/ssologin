package main

import (
	_ "github.com/doggytty/ssologin/routers"
	_ "github.com/astaxie/beego/session/mysql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/doggytty/ssologin/models"
	"fmt"
	"github.com/doggytty/ssologin/utils"
	"github.com/astaxie/beego/context"
)

func main() {
	// beego默认logger
	logs.Async()
	// log异步输出
	logs.EnableFuncCallDepth(true)
	if beego.BConfig.RunMode == "dev" {
		logs.SetLogger(logs.AdapterConsole)
	} else {
		logs.SetLogger(logs.AdapterFile, `{"filename":"main.log","daily":true,"maxdays":10}`)
	}

	// init database
	logs.GetBeeLogger().Debug("init database!")
	models.SyncDataBase()
	logs.GetBeeLogger().Debug("init database success!")

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

	// xsrf
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "61oETzYUQ62dkL5gEmGeJJFuYh7EQnp2XdTP1o"
	beego.BConfig.WebConfig.XSRFExpire = 3600  //过期时间，默认1小时

	// 启用session
	beego.BConfig.WebConfig.Session.SessionOn = true
	// 使用数据库存储session
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	beego.BConfig.WebConfig.Session.SessionProviderConfig = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
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
	//beego.InsertFilter("/*", beego.BeforeStatic, filters.FilterBlackDDOS)
	// 2、session filter
	//beego.InsertFilter("/*", beego.BeforeRouter, filters.FilterLogin)
	// 3、admin filter
	//beego.InsertFilter("/admin/*", beego.BeforeRouter, filters.FilterAdministrator)
	// 设置form表单支持put\delete
	var filterMethod = func(ctx *context.Context) {
		if ctx.Input.Query("_method")!="" && ctx.Input.IsPost(){
			ctx.Request.Method = ctx.Input.Query("_method")
		}
	}
	beego.InsertFilter("*", beego.BeforeRouter, filterMethod)

	// 自定义template函数
	beego.AddFuncMap("paginationJump",utils.Jump)
	beego.AddFuncMap("paginationPrefix",utils.Prefix)
	beego.AddFuncMap("paginationSuffix",utils.Suffix)
	beego.AddFuncMap("paginationShowPrefix",utils.ShowPrefix)
	beego.AddFuncMap("paginationShowSuffix",utils.ShowSuffix)
	beego.AddFuncMap("paginationShowFirst",utils.ShowFirst)
	beego.AddFuncMap("paginationShowLast",utils.ShowLast)
	beego.AddFuncMap("paginationGetPageNumber",utils.GetPageNumber)
	beego.AddFuncMap("paginationGetBeginIndex",utils.GetBeginIndex)
	beego.AddFuncMap("paginationGetEndIndex",utils.GetEndIndex)
	beego.AddFuncMap("paginationGetTotalPage",utils.GetTotalPage)

	// 启动beego
	beego.Run("127.0.0.1:8080")
}

