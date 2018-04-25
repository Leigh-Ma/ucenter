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

func (r *JResp) Success(d ...*D) *JResp{
	r.ErrorCode = OK
	r.ErrorReason = "success"
	if len(d) > 0 {
		r.Data = d[0]
	}
	return r
}

func (r *JResp) Error(code uint, more ...string) *JResp {
	r.ErrorCode = code
	r.ErrorReason = ErrDesc[code]

	if len(more) > 0 {
		r.ErrorReason += ": " + more[0]
	}
	return r
}

func (r *JResp) Set(code uint, d ...*D) *JResp {
	r.Error(code)
	if len(d) > 0 {
		r.Data = d[0]
	}
	return r
}
