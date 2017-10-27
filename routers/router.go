package routers

import (
	"github.com/doggytty/ssologin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.AdminController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/admin", &controllers.AdminController{})


    beego.Router("/noAuth", &controllers.LoginController{}, "get:NoAuth")
    beego.Router("/notAdmin", &controllers.LoginController{}, "get:NotAdmin")
}
