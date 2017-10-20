package models

// 子系统
type SubSystem struct {
	Sid string `orm:"size(20);pk;column(sid)"`
	SName string `orm:"size(50);column(sys_name)"`
	ClientId string `orm:"size(100);column(client_id)"`
	ClientSecret string `orm:"size(100);column(client_secret)"`
	SUrl string `orm:"size(100);column(sys_url)"`
	Status string `orm:"size(100);column(sys_status)"`
}
func (ui *SubSystem) TableName() string {
	return "sub_system"
}

// 系统属性表格
type SubProperties struct {
	Sid string `orm:"size(10);column(sid)"`
	Key string `orm:"size(100);column(prop_key)"`
	Value string `orm:"size(100);column(prop_value)"`
}
func (ui *SubProperties) TableName() string {
	return "sub_properties"
}


