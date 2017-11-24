package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"github.com/doggytty/goutils/encrypt"
	"errors"
	"fmt"
)

type Administrator struct {
	Email string `orm:"size(50);pk;column(email)"`
}

func (ad *Administrator) InsertAdministrator(email string) bool {
	o := orm.NewOrm()
	admin := new(Administrator)
	admin.Email = email
	num, err := o.Insert(admin)
	if err != nil {
		logger.Error("[Administrator]insert administrator failed!", err)
		return false
	}
	if num == 1 {
		logger.Debug("[Administrator]insert administrator success!", num)
		return true
	}
	logger.Error("[Administrator]insert administrator wrong!", num)
	return false
}


// 用户表格
type UserInfo struct {
	Uid string `orm:"pk;size(32);column(uid)"`	// 32位随机字符串
	Email string `orm:"size(50);column(email);unique"`
	Password string `orm:"size(50);column(password)"`
	NickName string `orm:"size(30);column(nickname)"`
	Age int `orm:"null;column(age)"`
	Sex int `orm:"null;column(sex);default(0)"`
	IdNo string `orm:"null;size(18);column(idno)"`

	Status int `orm:"column(status);default(0)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime);column(update_time)"`
}

func (ui *UserInfo) TableName() string {
	return "user_info"
}

func (ui *UserInfo) IsAdministrator(uid string) *UserInfo {
	logger.Debug("is user %s administrator?", ui.Email)
	o := orm.NewOrm()
	err := o.Raw("SELECT ui.* FROM user_info ui, administrator ad WHERE ui.email=ad.email and ui.uid = ?", uid).QueryRow(&ui)
	if err != nil {
		logger.Error("query administrator failed!", err)
		return nil
	}
	return ui
}

func (ui *UserInfo) GetById(uid string) *UserInfo {
	o := orm.NewOrm()
	ui.Uid = uid
	err := o.Read(&ui)
	if err == orm.ErrNoRows {
		logger.Error("query not exist:", uid)
	} else if err == orm.ErrMissPK {
		logger.Error("can not find pk", uid)
	} else {
		return ui
	}
	return nil
}

func (ui *UserInfo) SetPassword(password string) {
	md5password := encrypt.Md5String(password)
	ui.Password = md5password
}

func (ui *UserInfo) CheckUserInfo(email, password string) (*UserInfo, error) {
	ui.Email = email
	ui.SetPassword(password)

	o := orm.NewOrm()
	qs := o.QueryTable(ui)

	var user []*UserInfo
	num, err := qs.Filter("email", email).Filter("password", ui.Password).All(&user)
	if err != nil {
		logger.Error("query by username,password failed, %s", err)
		return nil, err
	}
	if num == 1 {
		logger.Debug("query success")
		return user[0], nil
	}
	logger.Error("more than one person email is %s", email)
	return nil, errors.New(fmt.Sprintf("there is %d user email is %s", num, email))
}

func (ui *UserInfo) LastUserInfo(beginIndex, pageSize int) []*UserInfo {
	o := orm.NewOrm()
	qs := o.QueryTable(ui)

	var users []*UserInfo
	querySetter := qs.OrderBy("-create_time")
	if beginIndex == 0 {
		querySetter = querySetter.Limit(pageSize)
	} else {
		querySetter = querySetter.Limit(beginIndex, pageSize)
	}
	num, err := querySetter.All(&users)
	if err != nil {
		logger.Error("query userInfo failed", err)
		return nil
	}
	logger.Debug("query userInfo %d", num)
	return users
}


// 用户登录/退出系统记录
type UserLogin struct {
	Id int `orm:"pk;auto;column(id)"`
	Uid string `orm:"size(32);column(uid);index(idx_uid)"`
	LoginIp string `orm:"column(login_ip)"`	// 登陆ip
	Sid string `orm:"column(sid);index(idx_sid)"`	// 登录的子系统
	IsLogin bool `orm:"column(is_login)"`	// 是否登录
	IsNormalExit int `orm:"column(is_normal_exit);default(0)"`	// 是否正常退出
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"`
}

func (ul *UserLogin) TableName() string {
	return "user_login"
}

func (ul *UserLogin) InsertUserLogin(uid, ip, appId string, isLogin bool) int {
	o := orm.NewOrm()
	ul.Uid = uid
	ul.LoginIp = ip
	ul.Sid = appId
	ul.IsLogin = isLogin
	ulId, err := o.Insert(ul)
	if err == nil {
		return int(ulId)
	}
	return -1
}

func (ul *UserLogin) LastUserLogin(beginIndex, pageSize int) []*UserLogin {
	o := orm.NewOrm()
	qs := o.QueryTable(ul)
	var userLogin []*UserLogin
	querySetter := qs.OrderBy("-id")
	if beginIndex == 0 {
		querySetter = querySetter.Limit(pageSize)
	} else {
		querySetter = querySetter.Limit(beginIndex, pageSize)
	}
	num, err := querySetter.All(&userLogin)
	if err != nil {
		logger.Error("query userLogin failed", err)
		return nil
	}
	logger.Debug("query userLogin %d", num)
	return userLogin
}


// 用户-系统对应关系
type UserAuth struct {
	Id int `orm:"pk;auto;column(id)"`
	Uid string `orm:"size(32);column(uid)"`
	Sid string `orm:"size(20);column(sid)"`
	IsAdministrator bool `orm:"column(is_administrator)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"`
}

func (ua *UserAuth) TableName() string {
	return "user_auth"
}

