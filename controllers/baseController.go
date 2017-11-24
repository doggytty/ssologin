package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

var logger = logs.NewLogger()

func init()  {
	logger.Async()
	logger.EnableFuncCallDepth(true)
	if beego.BConfig.RunMode == "dev" {
		logger.SetLogger(logs.AdapterConsole)
	} else {
		logger.SetLogger(logs.AdapterFile, `{"filename":"controller.log","daily":true,"maxdays":10}`)
	}
}


type BaseController struct {
	beego.Controller
}

func (base *BaseController) Init(ctx *context.Context, controllerName, actionName string, app interface{})  {
	// 调用默认初始化设置
	logger.Debug("base controller Init")
	base.Controller.Init(ctx, controllerName, actionName, app)
}


func (base *BaseController) Prepare() {
	logger.Debug("base controller Prepare")
	base.Layout = "layout/base.html"
	base.LayoutSections = make(map[string]string)
	base.LayoutSections["Header"] = "layout/header.html"
	base.LayoutSections["Scripts"] = "layout/script.html"
	base.LayoutSections["Styles"] = "layout/styles.html"
	base.LayoutSections["Footer"] = "layout/footer.html"
	base.LayoutSections["Navigation"] = "layout/navigation.html"
	// 判断用户是否登陆
	uid := base.GetSession("uid")
	base.Data["IsLogin"] = uid != nil


	//if ctx.BeegoInput.Query("_method")!="" && ctx.BeegoInput.IsPost(){
	//this.Ctx.Input.IsPost()
	//log.Println(this.Ctx.Input.IsPost())
	//处理uid,username

	//tmpId := base.Ctx.Input.CruSession.Get("uid")
	//if tmpId == nil {
	//	return
	//}
	//uid := tmpId.(int64)
	//if uid > 0 {
	//	if this.Ctx.Input.IsAjax() {
	//		log.Println("ajax request!")
	//	} else {
	//		user := m.GetUserById(uid)
	//		this.Data["uid"] = uid
	//		this.Data["username"] = user.Username
	//	}
	//}
	//log.Println(this.Ctx.Request)
}


func (base *BaseController) Rsp(result bool, err string) {
	base.Data["json"] = map[string]interface{}{"result": result, "msg": err}
	base.ServeJSON()
}

