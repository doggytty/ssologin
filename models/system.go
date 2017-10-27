package models

import "github.com/astaxie/beego/orm"

// 子系统
type SubSystem struct {
	Sid string `orm:"size(20);pk;column(sid)"`
	SName string `orm:"size(50);column(sys_name)"`
	ClientId string `orm:"size(100);column(client_id)"`	// clientid=md5(sid\md5(secret))
	ClientSecret string `orm:"size(100);column(client_secret)"`
	SUrl string `orm:"size(100);column(sys_url)"`
	Status string `orm:"size(100);column(sys_status)"`
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
		logger.Error("query not exist:", appId)
	} else if err == orm.ErrMissPK {
		logger.Error("can not find pk", appId)
	} else {
		return tt
	}
	return nil
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
