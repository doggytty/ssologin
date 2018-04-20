package routers

import (
	"github.com/doggytty/ssologin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.AdminController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/admin", &controllers.AdminController{})

    beego.Router("/status/:uid", &controllers.LoginController{}, "get:Status")
    beego.Router("/noAuth", &controllers.LoginController{}, "get:NoAuth")
    beego.Router("/notAdmin", &controllers.LoginController{}, "get:NotAdmin")

	beego.Router("/admin/index", &controllers.AdminController{})

	beego.Router("/admin/system/list", &controllers.SystemController{}, "get:List")
	beego.Router("/admin/system", &controllers.SystemController{})
	beego.Router("/admin/system/:sid", &controllers.SystemController{})
	beego.Router("/admin/system/delete/:sid", &controllers.SystemController{}, "get:DeleteSystem")

	beego.Router("/admin/userinfo", &controllers.AdminController{})
	beego.Router("/admin/sysuser", &controllers.AdminController{})
	beego.Router("/online/monitor", &controllers.AdminController{})
	beego.Router("/about", &controllers.AdminController{})
}
