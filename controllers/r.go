package controllers

import (
	"github.com/astaxie/beego"
	"strings"
)

type IExport interface {
	Export() func(string)
}

type RouterGroup struct {
	Namespace string
	Routers   map[string]IExport
}

func (c *RouterGroup) RegisterRouter() {
	for name, router := range c.Routers {
		router.Export()(c.Namespace + name) //in case of false call to c.Export
	}
}

func Exportor(ctrl beego.ControllerInterface, r map[string]string) func(string) {
	//"GET: /index" : "Index"

	return func(ctrlNameSpace string) {
		for route, fn := range r {
			ss := strings.SplitN(route, ":", 2)
			match := strings.Trim(ss[1], " ")
			method := strings.ToLower(strings.Trim(ss[0], " "))

			path := ctrlNameSpace + match
			call := method + ":" + fn

			beego.Info(path, " ", call)
			beego.Router(path, ctrl, call)
		}
	}
}
