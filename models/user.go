package models


import (
	"ucenter/library/types"
	"sync"
)

type User struct {
	TCom
	UserName string `orm:"unique;size(64)"`
	Email    string `orm:"unique;size(64)"`
	Password string
	Salt     string
	UUID     string
	IsActive bool
	IsForbid bool
	RawPwd   string
}

func (t *User) TableName() string {
	return "users"
}

func (t *User) md5Pwd(pwd string) string {
	if len(t.Salt) == 0 {
		t.Salt = types.RandomString(32)
	}

	return types.MD5(pwd + t.Salt)
}

func (t *User) SetPassword(pwd string) {
	t.Salt = ""
	t.Password = t.md5Pwd(pwd)
}

func (t *User) VerifyPassword(pwd string) bool{
	return types.MD5(pwd + t.Salt) == t.Password
}

var (
	UserM = NewUserManager()
)

type UserManager struct {
	sync.RWMutex
	UsersById   map[int64]*User
	UsersByUUID map[string]*User
	UsersByEmail map[string]*User
}

func NewUserManager() *UserManager{
	return &UserManager{
		UsersById: make(map[int64]*User, 0),
		UsersByUUID: make(map[string]*User, 0),
		UsersByEmail: make(map[string]*User, 0),
	}
}

func (t *UserManager) AddUser(user *User) {
	t.Lock()
	t.UsersById[user.Id] = user
	t.UsersByUUID[user.UUID] = user
	if user.Email != "" {
		t.UsersByEmail[user.Email] = user
	}
	t.Unlock()
}

func (t *UserManager) GetVisitor(uuid string)(*User, bool) {
	t.RLock()
	user, ok := t.UsersByUUID[uuid]
	t.Unlock()
	isNew := false
	if !ok {
		user := &User{
			UUID:     uuid,
			IsActive: true,
		}
		if err := user.FindBy("UUID", uuid, user); err != nil {
			user.Insert(user)
			isNew = true
		}

		t.AddUser(user)
	}

	return user, isNew
}

func (t *UserManager) GetByEmail(email string, uuid... string)(u *User, isNew bool) {

	t.RLock()
	user, ok := t.UsersByEmail[email]
	t.Unlock()

	if !ok {
		user := &User{
			Email:    email,
			IsActive: true,
		}

		if len(uuid) > 0 && len(uuid[0]) > 0 {
			user.UUID = uuid[0]
		}

		if err := user.FindBy("Email", email, user); err != nil {
			//user.Insert(user)
			isNew = true
		}

		t.AddUser(user)
	}


	u = user
	return
}