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

func (ui *SubSystem) CountByQuery(paramMap map[string]string) int {
	o := orm.NewOrm()
	qs := o.QueryTable(ui)
	if paramMap == nil || len(paramMap) == 0 {
		result, err := qs.Count()
		if err != nil {
			logger.Error("query subsystem count failed, %v", err)
			return 0
		}
		return int(result)
	}
	querySetter := qs.SetCond(orm.NewCondition())
	sid, ok := paramMap["sid"]
	if ok {
		querySetter = querySetter.Filter("sid__icontains", sid)
	}
	sysName, ok := paramMap["sys_name"]
	if ok {
		querySetter = querySetter.Filter("sys_name__icontains", sysName)
	}
	clientId, ok := paramMap["client_id"]
	if ok {
		querySetter = querySetter.Filter("client_id__icontains", clientId)
	}
	clientSecret, ok := paramMap["client_secret"]
	if ok {
		querySetter = querySetter.Filter("client_secret__icontains", clientSecret)
	}
	sysUrl, ok := paramMap["sys_url"]
	if ok {
		querySetter = querySetter.Filter("sys_url__icontains", sysUrl)
	}
	result, err := querySetter.Count()
	if err != nil {
		logger.Error("query subsystem count failed, %v", err)
		return 0
	}
	return int(result)
}

// 查询最新的10个系统
func (ui *SubSystem) QueryByPage(beginIndex, pageSize int, paramMap map[string]string) []*SubSystem {
	o := orm.NewOrm()
	qs := o.QueryTable(ui)

	var subSystem []*SubSystem
	querySetter := qs.OrderBy("-create_time")
	if beginIndex == 0 {
		querySetter = querySetter.Limit(pageSize)
	} else {
		querySetter = querySetter.Limit(beginIndex, pageSize)
	}
	if paramMap != nil && len(paramMap) > 0 {
		sid, ok := paramMap["sid"]
		if ok {
			querySetter = querySetter.Filter("sid__icontains", sid)
		}
		sysName, ok := paramMap["sys_name"]
		if ok {
			querySetter = querySetter.Filter("sys_name__icontains", sysName)
		}
		clientId, ok := paramMap["client_id"]
		if ok {
			querySetter = querySetter.Filter("client_id__icontains", clientId)
		}
		clientSecret, ok := paramMap["client_secret"]
		if ok {
			querySetter = querySetter.Filter("client_secret__icontains", clientSecret)
		}
		sysUrl, ok := paramMap["sys_url"]
		if ok {
			querySetter = querySetter.Filter("sys_url__icontains", sysUrl)
		}
	}
	num, err := querySetter.All(&subSystem)
	if err != nil {
		logger.Error("query subsystem failed", err)
		return nil
	}
	logger.Debug("query subsytem %d", num)
	return subSystem
}

func (ui *SubSystem) ModifySubSystem() bool {
	o := orm.NewOrm()
	// 获取 QuerySeter 对象，user 为表名
	qs := o.QueryTable("user")
	params := orm.Params{}
	if ui.SName != "" {
		params["sys_name"] = ui.SName
	}
	if ui.ClientId != "" {
		params["client_id"] = ui.ClientId
	}
	if ui.ClientSecret != "" {
		params["client_secret"] = ui.ClientSecret
	}
	if ui.SUrl != "" {
		params["sys_url"] = ui.SUrl
	}
	if ui.Status != "" {
		params["sys_status"] = ui.Status
	}
	num, err := qs.Filter("sid", ui.Sid).Update(params)
	if err == nil {
		return num == 1
	}
	return false
}

func (ui *SubSystem) DeleteSubSystem(sid string) bool {
	o := orm.NewOrm()
	ui.Sid = sid
	num, err := o.Delete(ui)
	if err == nil {
		return num == 1
	}
	return false
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
