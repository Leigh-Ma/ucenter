package main

import (
	_ "ucenter/routers"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/mysql"
	"github.com/astaxie/beego/orm"
	"os"
)

func main() {
	initOrm()
	beego.Run()
}


func initOrm() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/UCENTER?charset=utf8&loc=UTC")

	orm.RunCommand()
	orm.DebugLog = orm.NewLog(os.Stdout)

	orm.RunSyncdb("default", false, false)
}

