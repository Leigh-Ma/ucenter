package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/session/mysql"
	"os"
	"ucenter/library/types"
	_ "ucenter/routers"
)

func main() {
	initOrm()
	beego.Run()
}

func initOrm() {
	types.InitIDGen("123")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/UCENTER?charset=utf8&loc=UTC")

	orm.RunCommand()
	orm.DebugLog = orm.NewLog(os.Stdout)

	orm.RunSyncdb("default", false, false)
}
