package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Administrator struct {
	Email string `orm:"size(50);column(email);unique"`
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

func (ui *UserInfo) IsAdministrator() bool {
	logger.Debug("is user %s administrator?", ui.Email)
	o := orm.NewOrm()
	qs := o.QueryTable(new(Administrator))
	count, err := qs.Filter("email", ui.Email).Count()
	if err != nil {
		logger.Error("query administrator failed!", err)
		return false
	}
	if count == 1 {
		return true
	}
	return false
}


// 用户登录/退出系统记录
type UserLogin struct {
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

func (ul *UserLogin) InsertUserLogin() bool {

}




// 用户-系统对应关系
type UserAuth struct {
	Uid string `orm:"size(32);column(uid)"`
	Sid string `orm:"size(20);column(sid)"`
	IsAdministrator bool `orm:"column(is_administrator)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"`
}

func (ua *UserAuth) TableName() string {
	return "user_auth"
}

