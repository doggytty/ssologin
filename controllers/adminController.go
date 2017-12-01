package controllers

import (
	"github.com/doggytty/ssologin/models"
	"github.com/doggytty/goutils/htmlutil"
	"github.com/astaxie/beego"
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


// 第三方系统管理
type SystemController struct {
	AdminController
}

// 跳转到系统页面
func (sc *SystemController) Get() {
	// 分页查询系统列表
	pageSize, _ := sc.GetInt("pageSize", 15)
	pageIndex, _ := sc.GetInt("pageIndex", 1)
	// 其它查询条件
	sysIdString := sc.GetString("sid", "")
	sysNameString := sc.GetString("sys_name", "")
	clientIdString := sc.GetString("client_id", "")
	clientSecretString := sc.GetString("client_secret", "")
	sysUrlString := sc.GetString("sys_url", "")
	statusString := sc.GetString("status", "")
	timeString := sc.GetString("create_time", "")
	paramMap := make(map[string]string, 0)
	if sysIdString != "" {
		paramMap["sid"] = sysIdString
	}
	if sysNameString != "" {
		paramMap["sys_name"] = sysNameString
	}
	if clientIdString != "" {
		paramMap["client_id"] = clientIdString
	}
	if clientSecretString != "" {
		paramMap["client_secret"] = clientSecretString
	}
	if sysUrlString != "" {
		paramMap["sys_url"] = sysUrlString
	}
	if statusString != "" {
		paramMap["status"] = statusString
	}
	if timeString != "" {
		paramMap["create_time"] = timeString
	}

	beginIndex := pageSize * (pageIndex - 1)
	subSystem := new(models.SubSystem)
	pageSystem := subSystem.QueryByPage(beginIndex, pageSize, paramMap)
	totalRecord := subSystem.CountByQuery(paramMap)

	p := htmlutil.NewPagination(pageSize, pageIndex, totalRecord)
	sc.Data["PageSystem"] = pageSystem
	sc.Data["Pagination"] = p
	sc.TplName = "admin/system.html"
}

func (sc *SystemController) DetailSystem ()  {
	sid := sc.Ctx.Input.Param("sid")
	logger.Debug("now sid :%s", sid)
	if sid == "new" {
		logger.Debug("now create subSystem!")
	} else {
		subSystem := new(models.SubSystem)
		subSystem = subSystem.GetSubSystemById(sid)
		sc.Data["SubSystem"] = subSystem
	}
	sc.TplName = "admin/systemdetail.html"
}

func (sc *SystemController) ModifySystem ()  {
	sid := sc.GetString("sid")
	sysName := sc.GetString("sysName")
	clientId := sc.GetString("clientId")
	clientSecret := sc.GetString("clientSecret")
	sysUrl := sc.GetString("sysUrl")
	subSystem := new(models.SubSystem)
	subSystem.Sid = sid
	subSystem.SName = sysName
	subSystem.ClientId = clientId
	subSystem.ClientSecret = clientSecret
	subSystem.SUrl = sysUrl
	logger.Debug("now sid :%s", sid)
	result := subSystem.ModifySubSystem()

	if result {
		logger.Debug("now create subSystem!")
		sc.Redirect("/admin/system", 302)
	} else {
		flash := beego.NewFlash()
		flash.Set("status", "11011")
		flash.Store(&sc.Controller)
		sc.Redirect("/noAuth", 302)
	}
}

func (sc *SystemController) Delete ()  {
	sid := sc.GetString("sid")
	logger.Debug("now delete sid :%s", sid)
	subSystem := new(models.SubSystem)
	result := subSystem.DeleteSubSystem(sid)
	if result {
		logger.Debug("delete subSystem success!")
		sc.TplName = "admin/systemdetail.html"
	} else {
		logger.Debug("delete subSystem failed!")
	}

}

