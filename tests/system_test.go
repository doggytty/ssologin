package test

import (
	"path/filepath"
	"fmt"
	"github.com/astaxie/beego"
	"testing"
	"runtime"
	"github.com/doggytty/ssologin/models"
)

func preInit() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	apppath = filepath.Dir(apppath)
	fmt.Println(apppath)
	beego.TestBeegoInit(apppath)
}

func TestSubSystem_CountByQuery(t *testing.T) {
	models.SyncDataBase()
	subSystem := new(models.SubSystem)
	paramMap := make(map[string]string, 0)
	paramMap["sid"] = "z"
	count := subSystem.CountByQuery(paramMap)
	fmt.Println(count)
}