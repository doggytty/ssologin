package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"strconv"
	"time"
	"github.com/doggytty/goutils/encrypt"
	"fmt"
	"github.com/doggytty/ssologin/models"
	"html/template"
)

type LoginController struct {
	beego.Controller
}

func (l *LoginController) autoLogin()  {
	// 判断是否自动登录(放弃该功能,安全受限)
	userLoginCookie, ok := l.Ctx.GetSecureCookie("assistant", "userLoginCookie")
	if ok {
		// 检查是否uid，password，expireTime，isAutoLogin
		logger.Debug("userLoginCookie:", userLoginCookie)
		cookieString, err := encrypt.Base64Decode(userLoginCookie, true)
		if err == nil {
			cookieArray := strings.Split(cookieString, ",")
			if len(cookieArray) == 4 {
				password := cookieArray[1]
				uid, err := strconv.Atoi(cookieArray[0])
				if err != nil {
					logs.Error("uid parse to int failed!")
				}
				expireTime, err := strconv.ParseInt(cookieArray[2], 10,64)
				if err != nil {
					logs.Error("expireTime parse to int64 failed!")
				}
				isAutoLogin, err := strconv.ParseBool(cookieArray[3])
				if err != nil {
					logs.Error("isAutoLogin parse to bool failed!")
				}
				if err == nil {
					// 判断是否要自动登录 && 时间未过期
					if isAutoLogin && time.Now().Unix() < expireTime {
						// 根据uid查询用户信息
						//models.UserInfo{}
						logs.Debug(uid, password)
					}
				}
			} else {
				logs.Error("userLoginCookie must be split by , and len = 4")
			}
		} else {
			logs.Error("cannot decode userLoginCookie", err)
		}
	} else {
		logs.Debug("no userLoginCookie")
	}

}

func (l *LoginController) Get() {
	redirectUrl := l.GetString("redirect_url")
	appId := l.GetString("app_id")
	clientSecret := l.GetString("client_secret")
	// session中有数据
	uid := l.GetSession("uid")
	if uid != nil {
		// 校验客户端是否
		subSystem := new(models.SubSystem)
		subSystem = subSystem.GetSubSystemById(appId)
		if subSystem == nil {
			logger.Error("app_id is not right!")
			flash := beego.NewFlash()
			flash.Set("status", "11011")
			flash.Store(&l.Controller)
			l.Redirect("/noAuth", 302)
			return
		}
		// 校验app_id\client_secret
		client_id := encrypt.Md5String(fmt.Sprintf("%s:%s", appId, encrypt.Md5String(clientSecret)))
		if client_id != subSystem.ClientId {
			logger.Error("app_id/client_secret is not right!")
			flash := beego.NewFlash()
			flash.Set("status", "11012")
			flash.Store(&l.Controller)
			l.Redirect("/noAuth", 302)
			return
		}
		logger.Debug("user %s is online, no return to %s", uid, redirectUrl)
		// 写日志
		ul := new(models.UserLogin)
		ulId := ul.InsertUserLogin(uid.(string), l.Ctx.Input.IP(), appId, true)
		if ulId < 0 {
			logger.Error("user login into db failed!")
		}
		// 跳转到对应的url
		tmpUrl := fmt.Sprintf("%s?uid=%s&redirect_url=%s", subSystem.SUrl, uid, redirectUrl)
		logger.Debug(tmpUrl)
		l.Redirect(tmpUrl, 302)
	} else {
		l.Data["xsrfdata"] = template.HTML(l.XSRFFormHTML())
		l.Data["RedirectUrl"] = redirectUrl
		l.Data["AppId"] = appId
		l.Data["ClientSecret"] = clientSecret
		l.TplName = "login.html"
	}
}

func (l *LoginController) Post() {
	redirectUrl := l.GetString("redirect_url")
	appId := l.GetString("app_id")
	clientSecret := l.GetString("client_secret")
	email := l.GetString("email")
	password := l.GetString("password")

	// 校验客户端是否注册
	subSystem := new(models.SubSystem)
	subSystem = subSystem.GetSubSystemById(appId)
	if subSystem == nil {
		logger.Error("app_id is not right!")
		flash := beego.NewFlash()
		flash.Set("status", "11011")
		flash.Store(&l.Controller)
		l.Redirect("/noAuth", 302)
		return
	}
	// 校验app_id\client_secret
	clientId := encrypt.Md5String(fmt.Sprintf("%s:%s", appId, encrypt.Md5String(clientSecret)))
	if clientId != subSystem.ClientId {
		logger.Error("app_id/client_secret is not right!")
		flash := beego.NewFlash()
		flash.Set("status", "11012")
		flash.Store(&l.Controller)
		l.Redirect("/noAuth", 302)
		return
	}
	// 判断用户是否在共享session中有数据
	uid := l.GetSession("uid")
	if uid != nil {
		logger.Debug("user %s is online, no return to %s", uid, redirectUrl)
		// 写日志
		ul := new(models.UserLogin)
		ulId := ul.InsertUserLogin(uid.(string), l.Ctx.Input.IP(), appId, true)
		if ulId < 0 {
			logger.Error("user login into db failed!")
		}
	} else {
		// user login
		userInfo := new(models.UserInfo)
		userInfo, err := userInfo.CheckUserInfo(email, password)
		if err != nil {
			logger.Error("username password is not right!")
			flash := beego.NewFlash()
			flash.Set("status", "11014")
			flash.Store(&l.Controller)
			l.Redirect("/noAuth", 302)
			return
		}
		if userInfo.Status == -1 {
			logger.Error("user is probi password is not right!")
			flash := beego.NewFlash()
			flash.Set("status", "11014")
			flash.Store(&l.Controller)
			l.Redirect("/noAuth", 302)
			return
		}
		// set session
		uid = userInfo.Uid
		l.SetSession("uid", userInfo.Uid)
	}
	// todo 发送用户状态到用户状态信道
	tmpUrl := fmt.Sprintf("%s?uid=%s&redirect_url=%s", subSystem.SUrl, uid, redirectUrl)
	logger.Debug(tmpUrl)
	// 将用户uid添加到
	l.TplName = "index.html"
	l.Redirect(tmpUrl, 302)
}

func (l *LoginController) NoAuth()  {
	flash := beego.ReadFromRequest(&l.Controller)
	if status, ok:=flash.Data["status"]; ok{
		if l.IsAjax() {
			l.Rsp(false, status)
		} else {
			l.Data["status"] = status
			l.TplName = "noAuth.html"
		}
	} else {
		if l.IsAjax() {
			l.Rsp(false, "no auth to access!")
		} else {
			l.Data["status"] = "11010"	// 未知错误
			l.TplName = "noAuth.html"
		}
	}
}

func (l *LoginController) NotAdmin()  {
	if l.IsAjax() {
		l.Rsp(false, "user is not administrator!")
	} else {
		l.TplName = "notAdmin.tpl"
	}
}

func (l *LoginController) Status()  {
	uid := l.Ctx.Input.Param(":uid")
	// 判断用户是否在线
	// 检查用户的最后操作时间,使用redis
	l.Data["json"] = map[string]interface{}{"uid": uid, "status": "online"}
	l.ServeJSON()
	// l.Ctx.WriteString("hello")
}

func (l *LoginController) Rsp(result bool, err string) {
	l.Data["json"] = map[string]interface{}{"result": result, "msg": err}
	l.ServeJSON()
}

type LogoutController struct {
	beego.Controller
}

func (lc *LogoutController) Get()  {
	err := lc.CruSession.Flush()
	if err != nil {
		lc.DelSession("uid")
		lc.DestroySession()
	}
	lc.TplName = "login.html"
}


