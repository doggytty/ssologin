package filters

import (
	"github.com/astaxie/beego/context"
	"log"
	"github.com/doggytty/goutils/collections"
)

var ExcludeUrl = []string{"/login", "/logout"}

// 过滤是否登录
func FilterLogin(ctx *context.Context) {
	uid := ctx.Input.CruSession.Get("uid")
	currURI := ctx.Request.RequestURI
	// 检查是否必须过滤的url
	if collections.Contains(ExcludeUrl, currURI) {
		log.Printf("url: %s\n", currURI)
		ctx.Input.CruSession.Flush()
	} else {
		log.Printf("not url: %s\n", currURI)
		if uid == nil {
			ctx.Redirect(302, "/login")
		} else {
			uidValue, ok := uid.(int)
			if !ok {
				ctx.Redirect(302, "/login")
			} else {
				log.Printf("now uid %d\n", uidValue)
			}
		}
	}
}


// 过滤是否超级管理员
func FilterAdministrator(ctx *context.Context)  {

}




// 屏蔽黑名单访问,防止DDOS
func FilterBlackDDOS(ctx *context.Context)  {
	// 获取黑名单信息
	// 获取当前请求的ip地址是否在黑名单中(原始地址\)
}

