package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

// 子系统
type SubSystem struct {
	Sid string `orm:"size(20);pk;column(sid)"`
	SName string `orm:"size(50);column(sys_name)"`
	ClientId string `orm:"size(100);column(client_id)"`	// clientid=md5(sid\md5(secret))
	ClientSecret string `orm:"size(100);column(client_secret)"`
	SUrl string `orm:"size(100);column(sys_url)"`
	Status string `orm:"size(100);column(sys_status)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime);column(create_time)"`
}

func (ui *SubSystem) TableName() string {
	return "sub_system"
}

func (ui *SubSystem) GetSubSystemById(appId string) *SubSystem {
	o := orm.NewOrm()
	tt := new(SubSystem)
	tt.Sid = appId
	err := o.Read(tt)
	if err == orm.ErrNoRows {
		logger.Error("query not exist: %s", appId)
	} else if err == orm.ErrMissPK {
		logger.Error("can not find pk %s", appId)
	} else {
		return tt
	}
	return nil
}


// 查询最新的10个系统
func (ui *SubSystem) QueryByPage(beginIndex, pageSize int) []*SubSystem {
	o := orm.NewOrm()
	qs := o.QueryTable(ui)

	var subSystem []*SubSystem
	querySetter := qs.OrderBy("-create_time")
	if beginIndex == 0 {
		querySetter = querySetter.Limit(pageSize)
	} else {
		querySetter = querySetter.Limit(beginIndex, pageSize)
	}
	num, err := querySetter.All(&subSystem)
	if err != nil {
		logger.Error("query subsystem failed", err)
		return nil
	}
	logger.Debug("query subsytem %d", num)
	return subSystem
}


// 系统属性表格
type SubProperties struct {
	Id int `orm:"pk;auto;column(id)"`
	Sid string `orm:"size(10);column(sid)"`
	Key string `orm:"size(100);column(prop_key)"`
	Value string `orm:"size(100);column(prop_value)"`
}
func (ui *SubProperties) TableName() string {
	return "sub_properties"
}


//type Session struct {
//	SessionKey string `orm:"size(64);pk;column(session_key)"`
//	SessionData []byte `orm:"column(session_data)"`
//	SessionExpiry uint `orm:"column(session_expiry)"`
//}
//func (ui *Session) TableName() string {
//	return "session"
//}
