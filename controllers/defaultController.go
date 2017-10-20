package controllers

import (
	"github.com/astaxie/beego"
	"os"
)

type MainController struct {
	BaseController
}

func (m *MainController) Get() {
	// 展示系统的环境变量
	envKeys := os.Environ()
	envData := make(map[string]string, len(envKeys))
	for _, keyString := range envKeys {
		valueString := os.Getenv(keyString)
		envData[keyString] = valueString
	}
	m.Data["environ"] = envData

	m.TplName = "index.html"
	m.Data["NavType"] = "main"
}

















type DemoController struct {
	beego.Controller
}

func (c *DemoController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "demo.tpl"
}
