# Error code define

## Common
0-OK-Success

## Internal
1-INTERNAL_ERROR-Internal error
2-DATA_BASE_ERROR-Database operation error
3-PARAMS_ERROR-Input parameters error
4-WEB_SOCKET_NEEDED-Should enable web socket in http header
5-SERVE_JSON_ONLY-We just support json request type
6-PLEASE_RE_LOGIN-Please re login, user_id or token invalid
7-ORDER_PRE_CREATE_ERR-Create order on server error

## User
### Login
1001-EMAIL_NOT_REGISTERED-Email is not registered
1002-PASSWORD_ERROR-Password not correct
1003-TOKEN_EXPIRED-Token expired
1004-TOKEN_INVALID-Token invalid
1005-USER_ID_INVALID-User not found by user id given

1050-WX_AUTH_BY_CODE_ERR-Use WeChat code auth error

### User action
1101-PASSWORD_NOT_CHANGED-Password is not changed

### Register

### Sign tasks
1301-HAVE_SIGNED_TODAY-You have signed today
1302-HOUR_SIGN_LATER-You should gain hourly reward later

#### Password
2001-PASSWORD_MISMATCH-Password mismatch
2002-PASSWORD_INVALID-Password rules not satisfied
#### Email
2011-EMAIL_HAS_BEEN_TAKEN-Email has been taken
2012-EMAIL_INVALID-Email is not valid

## Word Battle
3001-WB_JOIN_BATTLE_FAILED-Join battle error


