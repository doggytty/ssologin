package controllers

import (
	"github.com/doggytty/goutils/htmlutil"
	"github.com/doggytty/goutils/regexputils"
	"github.com/doggytty/goutils/encrypt"
	"fmt"
	"time"
	"github.com/doggytty/goutils/stringutil"
	"github.com/doggytty/ssologin/models"
	"html/template"
)

// 第三方系统管理
type SystemController struct {
	AdminController
}

func (sc *SystemController) Prepare() {
	logger.Debug("admin controller Prepare")
	sc.AdminController.Prepare()
	sc.LayoutSections["Scripts"] = "layout/script_switch.html"
	sc.LayoutSections["Styles"] = "layout/styles_switch.html"
}

func (sc *SystemController) List() {
	// 分页查询系统列表
	pageSize, _ := sc.GetInt("pageSize", 15)
	pageIndex, _ := sc.GetInt("pageIndex", 1)
	paramMap := make(map[string]string)

	beginIndex := pageSize * (pageIndex - 1)
	subSystem := new(models.SubSystem)
	pageSystem := subSystem.QueryByPage(beginIndex, pageSize, paramMap)
	totalRecord := subSystem.CountByQuery(paramMap)

	p := htmlutil.NewPagination(pageSize, pageIndex, totalRecord)
	sc.Data["PageSystem"] = pageSystem
	sc.Data["Pagination"] = p
	sc.Data["xsrfdata"] = template.HTML(sc.XSRFFormHTML())
	sc.TplName = "admin/system.html"
}

// 跳转到系统页面
func (sc *SystemController) Get() {
	sc.Data["xsrfdata"] = template.HTML(sc.XSRFFormHTML())
	sc.TplName = "admin/systemdetail.html"
	// 判断是否编辑subsystem
	sid := sc.Ctx.Input.Param(":sid")
	if sid != "" {
		subSystem := new(models.SubSystem)
		subSystem = subSystem.GetSubSystemById(sid)
		sc.Data["SubSystem"] = subSystem
	}
}

// 新增系统
func (sc *SystemController) Post() {
	// 其它查询条件
	sysName := sc.GetString("sysName")
	sysUrl := sc.GetString("sysUrl")
	status := sc.GetString("status")
	subSystem := new(models.SubSystem)
	subSystem.SName = sysName
	subSystem.SUrl = sysUrl
	subSystem.Status = status == "1"
	sc.Data["SubSystem"] = subSystem
	sc.Data["xsrfdata"] = template.HTML(sc.XSRFFormHTML())
	// 检查url格式是否正确
	urlFlag := regexputils.IsStrictUrl(sysUrl)
	if !urlFlag {
		sc.Data["errorCode"] = "SE0001"// url格式错误
		sc.TplName = "admin/systemdetail.html"
		return
	}
	// 检查sysname格式是否正确
	nameFlag := regexputils.IsWordString(sysName)
	if !nameFlag {
		sc.Data["errorCode"] = "SE0002"// name格式错误
		sc.TplName = "admin/systemdetail.html"
		return
	}
	// 检查sysName是否存在

	nameCount := subSystem.CountByName(sysName)
	if nameCount < 0 {
		sc.Data["errorCode"] = "SE0003"// 数据库名称查询错误
		sc.TplName = "admin/systemdetail.html"
	} else if nameCount > 1 {
		sc.Data["errorCode"] = "SE0004"// 系统已存在
		sc.TplName = "admin/systemdetail.html"
	} else {
		sid := stringutil.GenerateStringsSize(16)
		clientSecret := stringutil.GenerateStringsSize(32)
		clientId := encrypt.Md5String(fmt.Sprintf("%s:%s", sid, encrypt.Md5String(clientSecret)))
		// 添加到数据库
		subSystem.ClientId = clientId
		subSystem.ClientSecret = clientSecret
		subSystem.Sid = sid
		subSystem.SName = sysName
		subSystem.Status = status == "1"
		subSystem.SUrl = sysUrl
		subSystem.CreateTime = time.Now()
		systemFlag := subSystem.InsertSubSystem()
		if !systemFlag {
			// insert failed
			sc.Data["errorCode"] = "SE0005"// 系统insert已存在
			sc.TplName = "admin/systemdetail.html"
		} else {
			sc.Redirect("/admin/system", 302)
		}
	}
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
	subSystem.SUrl = sysUrl
	subSystem.ClientId = clientId
	subSystem.ClientSecret = clientSecret
	sc.Data["SubSystem"] = subSystem
	sc.Data["xsrfdata"] = template.HTML(sc.XSRFFormHTML())
	// 检查url格式是否正确
	urlFlag := regexputils.IsStrictUrl(sysUrl)
	if !urlFlag {
		sc.Data["errorCode"] = "SE0001"// url格式错误
		sc.TplName = "admin/systemdetail.html"
		return
	}
	// 检查sysname格式是否正确
	nameFlag := regexputils.IsWordString(sysName)
	if !nameFlag {
		sc.Data["errorCode"] = "SE0002"// name格式错误
		sc.TplName = "admin/systemdetail.html"
		return
	}
	// 检查sysName是否存在

	nameCount := subSystem.CountByNameFilterSid(sysName, sid)
	if nameCount < 0 {
		sc.Data["errorCode"] = "SE0003"// 数据库名称查询错误
		sc.TplName = "admin/systemdetail.html"
	} else if nameCount > 1 {
		sc.Data["errorCode"] = "SE0004"// 系统已存在
		sc.TplName = "admin/systemdetail.html"
	} else {
		subSystem.Sid = sid
		subSystem.SName = sysName
		subSystem.ClientId = clientId
		subSystem.ClientSecret = clientSecret
		subSystem.SUrl = sysUrl
		modifyFlag := subSystem.ModifySubSystem()
		if !modifyFlag {
			// insert failed
			sc.Data["errorCode"] = "SE0006"// 系统update错误
			sc.TplName = "admin/systemdetail.html"
		} else {
			sc.Redirect("/admin/system", 302)
		}
	}
}

func (sc *SystemController) DeleteSystem ()  {
	sid := sc.Ctx.Input.Param(":sid")
	logger.Debug("now delete sid :%s", sid)
	subSystem := new(models.SubSystem)
	result := subSystem.DeleteSubSystem(sid)
	if result {
		logger.Debug("delete subSystem success!")
		// todo 删除user-system关系
		sc.Redirect("/admin/system", 302)
	} else {
		logger.Debug("delete subSystem failed!")
		subSystem := new(models.SubSystem)
		subSystem = subSystem.GetSubSystemById(sid)
		sc.Data["SubSystem"] = subSystem
		sc.Data["errorCode"] = "SE0007"// 系统delete错误
		sc.Data["xsrfdata"] = template.HTML(sc.XSRFFormHTML())
		sc.TplName = "admin/systemdetail.html"
	}
}
