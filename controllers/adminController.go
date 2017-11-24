package controllers

import (
	"github.com/doggytty/ssologin/models"
	"github.com/astaxie/beego/utils/pagination"
	"github.com/doggytty/goutils/htmlutil"
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
	lastSubSystem := subSystem.QueryByPage(0, 10)

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


// 第三方系统管理
type SystemController struct {
	AdminController
}

// 跳转到系统页面
func (sc *SystemController) Get() {
	// 分页查询系统列表
	pageSize, _ := sc.GetInt("pageSize", 15)
	pageIndex, _ := sc.GetInt("pageIndex", 1)

	p := new(htmlutil.Pagination)
	p.CurrentPage = pageIndex
	p.PageSize = pageSize



	beginIndex := pageSize * (pageIndex - 1)
	subSystem := new(models.SubSystem)
	pageSystem := subSystem.QueryByPage(beginIndex, pageSize)
	sc.Data["PageSystem"] = pageSystem
	sc.Data["CurrPage"] = pageIndex
	sc.Data["TotalPage"] = pageIndex
	// todo 在线人数模块
	sc.TplName = "admin/system.html"
}

func (sc *SystemController) ModifySystem ()  {
	sid := sc.GetString("sysId")
	logger.Debug("now sid :%s", sid)
}

func (sc *SystemController) Delete ()  {
	// 删除system
}

