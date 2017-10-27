package filters

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/context"
	"github.com/doggytty/ssologin/models"
	"strings"
	"github.com/doggytty/goutils/collections"
)

var ExcludeUrl = []string{"/login", "/logout", "/noAuth", "/notAdmin"}
var logger = logs.NewLogger()

func init()  {
	logger.Async()
	// log异步输出
	logger.EnableFuncCallDepth(true)
	if beego.BConfig.RunMode == "dev" {
		logger.SetLogger(logs.AdapterConsole)
	} else {
		logger.SetLogger(logs.AdapterFile, `{"filename":"filter.log","daily":true,"maxdays":10}`)
	}
}

// 过滤是否登录
func FilterLogin(ctx *context.Context) {
	uid := ctx.Input.CruSession.Get("uid")
	currURI := ctx.Request.RequestURI
	if strings.HasPrefix(currURI, "/login") {
		// 如果是登录,清理所有的session信息
		ctx.Input.CruSession.Flush()
		// todo 用户所有的监控状态重置
		logger.Debug("%s login system", uid)
	} else if collections.Contains(ExcludeUrl, currURI) {
		// 系统退出
		// todo 重置所有监控状态
		logger.Debug("%s logout system", uid)
	} else {
		logger.Debug("current url: %s", currURI)
		logger.Debug("current uid: %s", uid)
		if uid == nil {
			ctx.Redirect(302, "/login")
		} else {
			uidValue, ok := uid.(string)
			if !ok {
				ctx.Redirect(302, "/login")
			} else {
				logger.Debug("now uid %d\n", uidValue)
				// 向监控线程发送消息,用户在线状态
			}
		}
	}
}


// 过滤是否超级管理员
func FilterAdministrator(ctx *context.Context)  {
	uid := ctx.Input.CruSession.Get("uid")
	uidValue, ok := uid.(string)
	if !ok {
		ctx.Redirect(302, "/login")
		return
	}
	// 判断用户是否超级管理员
	ui := new(models.UserInfo)
	ui = ui.IsAdministrator(uidValue)
	if ui != nil {
		ctx.Redirect(302, "/notAdmin")
	}
}

// 屏蔽黑名单访问,防止DDOS
func FilterBlackDDOS(ctx *context.Context)  {
	// 获取黑名单信息
	// 获取当前请求的ip地址是否在黑名单中(原始地址\)
}

