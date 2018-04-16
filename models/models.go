package models

import "github.com/astaxie/beego/orm"

var (
	models = &MM{
		Tables: []ITable{
			&AuthToken{}, &User{}, &Player{},
		},
	}
)

func init() {
	models.registerModels()
}

type MM struct {
	Tables []ITable
}

func (m *MM) registerModels() {
	for _, table := range m.Tables {
		orm.RegisterModel(table)
	}
}
