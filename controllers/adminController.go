package controllers

//import (
//	"github.com/astaxie/beego"
//	"os"
//)

type AdminController struct {
	BaseController
}

func (m *AdminController) Get() {
	// 跳转到管理员页面, 仪表盘
	m.TplName = "index.html"
	m.Data["NavType"] = "main"
}
