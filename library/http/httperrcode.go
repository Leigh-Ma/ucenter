package http

const (
	OK = 0
	// Error code define
	//# Common
	ERR_OK = 0
	//# Internal
	ERR_INTERNAL_ERROR    = 1
	ERR_DATA_BASE_ERROR   = 2
	ERR_PARAMS_ERROR      = 3
	ERR_WEB_SOCKET_NEEDED = 4
	//# User
	//## Login
	ERR_EMAIL_NOT_REGISTERED = 1001
	ERR_PASSWORD_ERROR       = 1002
	ERR_TOKEN_EXPIRED        = 1003
	ERR_TOKEN_INVALID        = 1004
	ERR_USER_ID_INVALID      = 1005
	//## User action
	ERR_PASSWORD_NOT_CHANGED = 1101
	//## Register
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
	ERR_INTERNAL_ERROR:    "Internal error",
	ERR_DATA_BASE_ERROR:   "Database operation error",
	ERR_PARAMS_ERROR:      "Input parameters error",
	ERR_WEB_SOCKET_NEEDED: "Should enable web socket in http header",
	//# User
	//## Login
	ERR_EMAIL_NOT_REGISTERED: "Email is not registered",
	ERR_PASSWORD_ERROR:       "Password not correct",
	ERR_TOKEN_EXPIRED:        "Token expired",
	ERR_TOKEN_INVALID:        "Token invalid",
	ERR_USER_ID_INVALID:      "User not found by user id given",
	//## User action
	ERR_PASSWORD_NOT_CHANGED: "Password is not changed",
	//## Register
	//### Password
	ERR_PASSWORD_MISMATCH: "Password mismatch",
	ERR_PASSWORD_INVALID:  "Password rules not satisfied",
	//### Email
	ERR_EMAIL_HAS_BEEN_TAKEN: "Email has been taken",
	ERR_EMAIL_INVALID:        "Email is not valid",
	//# Word Battle
	ERR_WB_JOIN_BATTLE_FAILED: "Join battle error",
}
