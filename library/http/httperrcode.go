package http

const (
	OK = 0
	// Error code define
	//# Common
	ERR_OK = 0
	//# Internal
	ERR_INTERNAL_ERROR       = 1
	ERR_DATA_BASE_ERROR      = 2
	ERR_PARAMS_ERROR         = 3
	ERR_WEB_SOCKET_NEEDED    = 4
	ERR_SERVE_JSON_ONLY      = 5
	ERR_PLEASE_RE_LOGIN      = 6
	ERR_ORDER_PRE_CREATE_ERR = 7
	//# User
	//## Login
	ERR_EMAIL_NOT_REGISTERED = 1001
	ERR_PASSWORD_ERROR       = 1002
	ERR_TOKEN_EXPIRED        = 1003
	ERR_TOKEN_INVALID        = 1004
	ERR_USER_ID_INVALID      = 1005
	ERR_WX_AUTH_BY_CODE_ERR  = 1050
	//## User action
	ERR_PASSWORD_NOT_CHANGED = 1101
	//## Register
	//## Sign tasks
	ERR_HAVE_SIGNED_TODAY = 1301
	ERR_HOUR_SIGN_LATER   = 1302
	//### Password
	ERR_PASSWORD_MISMATCH = 2001
	ERR_PASSWORD_INVALID  = 2002
	//### Email
	ERR_EMAIL_HAS_BEEN_TAKEN = 2011
	ERR_EMAIL_INVALID        = 2012
	//# Word Battle
	ERR_WB_JOIN_BATTLE_FAILED = 3001
)

var ErrDesc = map[uint]string{
	// Error code define
	//# Common
	ERR_OK: "Success",
	//# Internal
	ERR_INTERNAL_ERROR:       "Internal error",
	ERR_DATA_BASE_ERROR:      "Database operation error",
	ERR_PARAMS_ERROR:         "Input parameters error",
	ERR_WEB_SOCKET_NEEDED:    "Should enable web socket in http header",
	ERR_SERVE_JSON_ONLY:      "We just support json request type",
	ERR_PLEASE_RE_LOGIN:      "Please re login, user_id or token invalid",
	ERR_ORDER_PRE_CREATE_ERR: "Create order on server error",
	//# User
	//## Login
	ERR_EMAIL_NOT_REGISTERED: "Email is not registered",
	ERR_PASSWORD_ERROR:       "Password not correct",
	ERR_TOKEN_EXPIRED:        "Token expired",
	ERR_TOKEN_INVALID:        "Token invalid",
	ERR_USER_ID_INVALID:      "User not found by user id given",
	ERR_WX_AUTH_BY_CODE_ERR:  "Use WeChat code auth error",
	//## User action
	ERR_PASSWORD_NOT_CHANGED: "Password is not changed",
	//## Register
	//## Sign tasks
	ERR_HAVE_SIGNED_TODAY: "You have signed today",
	ERR_HOUR_SIGN_LATER:   "You should gain hourly reward later",
	//### Password
	ERR_PASSWORD_MISMATCH: "Password mismatch",
	ERR_PASSWORD_INVALID:  "Password rules not satisfied",
	//### Email
	ERR_EMAIL_HAS_BEEN_TAKEN: "Email has been taken",
	ERR_EMAIL_INVALID:        "Email is not valid",
	//# Word Battle
	ERR_WB_JOIN_BATTLE_FAILED: "Join battle error",
}
