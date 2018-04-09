package http

type JResp struct {
	ErrorCode   uint        `json:"code"`
	ErrorReason string      `json:"msg,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

type D map[string]interface{}

func (d *D) Map() map[string]interface{} {
	return map[string]interface{}(*d)
}

func (d *D) Set(key string, value interface{}) {
	m := d.Map()
	m[key] = value
}

func (r *JResp) Success(d... *D) {
	r.ErrorCode = OK
	r.ErrorReason = "success"
	if len(d) > 0 {
		r.Data = d[0]
	}
}

func (r *JResp) ParamError(reason... string) {
	r.ErrorCode = ERR_PARAM_ERR
	if len(reason) > 0 {
		r.ErrorReason = reason[0]
	}
}

func (r *JResp) Error(reason... string) {
	r.ErrorCode = ERR_NORMAL_ERR
	if len(reason) > 0 {
		r.ErrorReason = reason[0]
	}
}

func (r *JResp) PasswordError(reason... string) {
	r.ErrorCode = ERR_PASSWORD_ERR
	if len(reason) > 0 {
		r.ErrorReason = reason[0]
	} else {
		r.ErrorReason = "Password Error"
	}
}