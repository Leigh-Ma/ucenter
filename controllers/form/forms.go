package form

import (
	"fmt"
	"net/url"
	"reflect"
	"ucenter/library/tools"
)

type FModifyPassword struct {
	Email       string
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
	PhoneID    string
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
	Uuid     string
	Email    string
	Password string
}

type FVisitorLogin struct {
	Uuid   string
	AppKey string
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
	fmt.Printf("%v\n", values)
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

		value := ""
		var vs []string
		if v, ok := values[tf.Name]; ok {
			vs = v
			if len(vs) > 0 {
				value = vs[0]
			}
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
