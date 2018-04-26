package form

import (
	"fmt"
	"net/url"
	"reflect"
	"ucenter/library/tools"
)

type FModifyPassword struct {
	Email       string `valid:"Email"`
	Password    string
	PasswordNew string
}

type FRegister struct {
	Name       string
	Email      string
	Password   string
	PasswordRe string
	Uuid       string
}

type FPhoneRegister struct {
	PhoneID    string `valid:"Phone"`
	VerifyCode string
}

type FTokenVerify struct {
	Token  string
	UserId int64
}

type FTokenLogin struct {
	UserId int64
	Token  string
}

type FPasswordLogin struct {
	Email    string `valid:"Email" json:"email"`
	Password string
}

type FVisitorLogin struct {
	Uuid   string
	AppKey string
}

type FSetPlayerName struct {
	Name string `valid:""`
}

type FBuyProduct struct {
	ProductId string
	Amount int //>=1
	Price float32
}

func shouldBeStructPtr(val reflect.Value) {
	if val.Kind() != reflect.Ptr && val.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("Try ParseForm form of no pointor to struct %s", val.Type().Name()))
	}

	if val.IsNil() {
		panic("Parse form cannot contain nil field")
	}
}

func ParseForm(form interface{}, values url.Values) {
	var (
		value string = ""
		ok  bool  = false
		vs []string = nil
	)

	valPtr := reflect.ValueOf(form)

	shouldBeStructPtr(valPtr)

	val := valPtr.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		vf := val.Field(i)
		tf := typ.Field(i)

		if !vf.CanSet() {
			//unexported  field
			continue
		}

		if tf.Anonymous && vf.CanAddr() {
			//Anonymous struct field
			ParseForm(vf.Addr().Interface(), values)
			return
		}

		tj := tf.Tag.Get("json")

		if vs, ok = values[tf.Name]; !ok {
			if tj != "" || tj != "-" {
				vs, ok = values[tj]
			}
		}

		value = ""
		if len(vs) > 0 {
			value = vs[0]
		}

		switch tf.Type.Kind() {
		case reflect.Bool:
			b, _ := tools.StrTo(value).Bool()
			vf.SetBool(b)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			x, _ := tools.StrTo(value).Int64()
			vf.SetInt(x)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			x, _ := tools.StrTo(value).Uint64()
			vf.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, _ := tools.StrTo(value).Float64()
			vf.SetFloat(x)
		case reflect.Struct:
			//not supported
		case reflect.String:
			vf.SetString(value)
		case reflect.Slice:
			vf.Set(reflect.ValueOf(vs))
		}
	}
}
