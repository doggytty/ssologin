package controllers

import (
	"github.com/doggytty/ssologin/models"
)

type AdminController struct {
	BaseController
}

func (m *AdminController) Prepare() {
	logger.Debug("admin controller Prepare")
	m.BaseController.Prepare()
	m.Data["NavType"] = "admin"
}

func (m *AdminController) Get() {
	// 跳转到管理员页面, 仪表盘
	subSystem := new(models.SubSystem)
	lastSubSystem := subSystem.QueryByPage(0, 10, nil)

	userInfo := new(models.UserInfo)
	lastUserInfo := userInfo.LastUserInfo(0, 10)

	userLogin := new(models.UserLogin)
	lastUserLogin := userLogin.LastUserLogin(0, 10)

	m.Data["lastSubSystem"] = lastSubSystem
	m.Data["lastUserInfo"] = lastUserInfo
	m.Data["lastUserLogin"] = lastUserLogin

	// todo 在线人数模块
	m.TplName = "admin/index.html"
}


