// pkg/apiresponse/messages.go
package apiresponse

var CodeMessages = []string{
	/* --------------------- 通用模块 --------------------- */
	"success",        // CodeOK
	"error",          // CodeErr
	"系统内部错误",   // CodeInternalServerError
	"请求参数不合法", // CodeInvalidRequestParams
	"未授权访问",     // CodeUnauthorizedAccess
	"无权限操作",     // CodeAccessDenied

	/* --------------------- 用户模块 --------------------- */
	"登录成功",                 // CodeUserLoginSuccess
	"登录失败（账号或密码错误）", // CodeUserLoginFailed
	"用户参数错误",             // CodeUserLoginParamError
	"登录密码未提供",           // CodeUserLoginPasswordMissing
	"登录用户名或邮箱未提供",   // CodeUserLoginIdentifierMissing
	"用户不存在",               // CodeUserNotFound
	"邮箱未注册",               // CodeUserEmailNotFound
	"用户账号被禁用",           // CodeUserDisabled
	"用户账号被锁定",           // CodeUserLocked
	"用户输入参数无效",         // CodeUserInputInvalid
	"用户未完成验证",           // CodeUserNotVerified
	"登录尝试次数过多",         // CodeUserLoginAttemptExceeded
	"用户会话已过期",           // CodeUserSessionExpired
	"密码哈希比对失败",         // CodeUserPasswordHashError
	"注册成功",                 // CodeUserRegisterSuccess
	"用户名已存在",             // CodeUserUsernameExists
	"邮箱已注册",               // CodeUserEmailExists
	"注册失败",                 // CodeUserRegisterFailed
	"验证码无效",               // CodeUserVerificationCodeInvalid
	"验证码已过期",             // CodeUserVerificationCodeExpired
	"信息更新成功",             // CodeUserProfileUpdateSuccess
	"信息更新失败",             // CodeUserProfileUpdateFailed
	"用户资料未完善",           // CodeUserProfileIncomplete
	"昵称长度超出限制",         // CodeUserNicknameTooLong
	"邮箱格式非法",             // CodeUserEmailFormatInvalid
	"手机号格式非法",           // CodeUserPhoneFormatInvalid

	/* --------------------- 用户 password --------------------- */
	"密码错误",         // CodeUserPasswordIncorrect
	"密码校验错误",     // CodeUserPasswordVerifyError
	"密码校验失败",     // CodeUserPasswordVerificationFailed
	"密码哈希比对失败", // CodeUserPasswordHashCompareFailed
	"密码过于简单",     // CodeUserPasswordTooWeak
	"密码格式非法",     // CodeUserPasswordFormatInvalid
	"原密码错误",       // CodeUserOldPasswordIncorrect

	/* --------------------- 用户 Token --------------------- */
	"Token 生成失败",         // CodeTokenGenerationFailed
	"Refresh Token 生成失败", // CodeRefreshTokenGenerationFailed
	"AccessToken 已过期",     // CodeAccessTokenExpired
	"RefreshToken 已过期",    // CodeRefreshTokenExpired
	"Token 无效",             // CodeTokenInvalidFormat
	"Refresh Token 不匹配",   // CodeRefreshTokenMismatch
	"Access Token 存储失败",  // CodeAccessTokenStorageFailed
	"Refresh Token 存储失败", // CodeRefreshTokenStorageFailed

	/* --------------------- 用户 Redis 其他 --------------------- */
	"用户设备绑定失败", // CodeUserRedisDeviceBindingFailed

	/* --------------------- Redis 模块 --------------------- */
	"Redis 连接失败",     // CodeRedisConnectionFailed
	"Redis 读取失败",     // CodeRedisReadFailed
	"Redis 写入失败",     // CodeRedisWriteFailed
	"Redis 删除失败",     // CodeRedisDeleteFailed
	"Redis 序列化失败",   // CodeRedisDataSerializationFailed ✅ 新增
	"Redis 反序列化失败", // CodeRedisDataDeserializationFailed ✅ 新增

	/* --------------------- 数据库模块 --------------------- */
	"数据库连接失败", // CodeDatabaseConnectFailed
	"数据库查询失败", // CodeDatabaseQueryFailed
	"数据库插入失败", // CodeDatabaseInsertFailed
	"数据库更新失败", // CodeDatabaseUpdateFailed
	"数据库删除失败", // CodeDatabaseDeleteFailed
}

func GetMessageByCode(code BaseCode) string {
	idx := int(code - BaseInitCode)
	if idx >= 0 && idx < len(CodeMessages) {
		return CodeMessages[idx]
	}
	return "未知错误"
}
