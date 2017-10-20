package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"strconv"
	"time"
	"github.com/doggytty/goutils/encrypt"
)

type LoginController struct {
	beego.Controller
}

func (l *LoginController) Get() {
	// 判断是否自动登录
	userLoginCookie, ok := l.Ctx.GetSecureCookie("assistant", "userLoginCookie")
	if ok {
		// 检查是否uid，password，expireTime，isAutoLogin
		logs.Debug("userLoginCookie:", userLoginCookie)
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


	l.TplName = "login.html"
}

func (l *LoginController) Post() {
	l.Data["Website"] = "beego.me"
	l.Data["Email"] = "astaxie@gmail.com"
	l.SetSession("uid", 110)
	//l.Ctx.Request.AddCookie()
	l.Redirect("/index", 302)
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
