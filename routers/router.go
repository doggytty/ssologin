package routers

import (
	"github.com/doggytty/ssologin/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
