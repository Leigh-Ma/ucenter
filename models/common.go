package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type iCom interface {
	SetId(id int64)
	GetId() int64
	IsNew() bool
	MarkOld() bool
}

type ITable interface {
	TableName() string
	IDB() *TCom
}

type TCom struct {
	Id        int64     `orm:"auto"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
	isNew     bool      `orm:"-"`
	dbh       *dbh      `orm:"-"`
}

func (t *TCom) IsNew() bool{
	return t.isNew
}

func (t *TCom) MarkOld() bool{
	i := t.isNew
	t.isNew = false
	return i
}

func (t *TCom) IDB() *TCom {
	return t
}

func (t *TCom) SetId(id int64) {
	t.Id = id
}

func (t *TCom) GetId() int64 {
	return t.Id
}

func (t *TCom) UseDBH(h *dbh) {
	t.dbh = h
}

func (t *TCom) Transaction(dbOperations func(*dbh) error) error {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.Transaction(dbOperations)
}

func (t *TCom) Reload(obj ITable) error {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.FindBy("id", t.Id, obj)
}

func (t *TCom) FindBy(field string, value interface{}, obj ITable) error {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.FindBy(field, value, obj)
}

func (t *TCom) FindById(id int64, obj ITable) error {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.FindBy("id", id, obj)
}

func (t *TCom) MultiQuery(cond *orm.Condition, table interface{}, cols ...string) ([]orm.Params, int64, error) {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.MultiQuery(cond, table, cols...)
}

func (t *TCom) NewQuery(obj interface{}) orm.QuerySeter {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.NewQuery(obj)
}

func (t *TCom) Insert(obj interface{}) (int64, error) {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.Insert(obj)
}

func (t *TCom) Update(obj interface{}, cols ...string) (int64, error) {
	if t.dbh == nil {
		t.dbh = DBH()
	}
	return t.dbh.Update(obj, cols...)
}

type dbh struct {
	orm.Ormer
}

func DBH() *dbh {
	return &dbh{
		Ormer: orm.NewOrm(),
	}
}

func (h *dbh) Transaction(dbOperations func(*dbh) error) error {
	h.Begin()
	err := dbOperations(h)
	if err != nil {
		h.Rollback()
		return err
	}
	return h.Commit()
}

func (h *dbh) FindBy(field string, value interface{}, obj ITable) error {
	return h.QueryTable(obj.TableName()).Filter(field, value).One(obj)
}

func (h *dbh) FindById(obj ITable, id ...int64) error {
	if len(id) > 0 {
		obj.IDB().SetId(id[0])
	}
	return h.Read(obj)
}

func (h *dbh) MultiQuery(cond *orm.Condition, table interface{}, cols ...string) ([]orm.Params, int64, error) {

	query := h.QueryTable(table).SetCond(cond)
	var container []orm.Params
	count, err := query.Values(&container, cols...)
	return container, count, err

}

func (h *dbh) NewQuery(obj interface{}) orm.QuerySeter {
	return h.QueryTable(obj)
}
