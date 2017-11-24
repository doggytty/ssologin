package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"
	_ "github.com/doggytty/ssologin/routers"

	"github.com/astaxie/beego"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/doggytty/goutils/stringutil"
	"fmt"
	"github.com/doggytty/goutils/encrypt"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	apppath = filepath.Dir(apppath)
	fmt.Println(apppath)
	beego.TestBeegoInit(apppath)
}


// TestBeego is a sample to run an endpoint test
func TestBeego(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	beego.Trace("testing", "TestBeego", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

func TestRandom(t *testing.T)  {
	app_id := stringutil.GenerateStringsSize(12)
	client_secret := stringutil.GenerateStringsSize(16)
	client_id := encrypt.Md5String(fmt.Sprintf("%s:%s", app_id, encrypt.Md5String(client_secret)))

	fmt.Println(fmt.Sprintf("app_id: %s", app_id))
	fmt.Println(fmt.Sprintf("client_secret: %s", client_secret))
	fmt.Println(fmt.Sprintf("client_id: %s", client_id))

	email := "sunlichuan@we.com"
	password := encrypt.Md5String("sunlichuan@we.com")
	fmt.Println(fmt.Sprintf("email: %s", email))
	fmt.Println(fmt.Sprintf("password: %s", password))
}

