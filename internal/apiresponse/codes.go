package apiresponse

// BaseCode 为整型状态码基准类型
type BaseCode int64

func (code BaseCode) String() BaseCode {
	return code // 输出如 "u10001"
}

// 响应码前缀（用于前端区分模块，例如："u10001"）
const (
	UserCodePrefix   = "u" // 用户模块响应码前缀
	CommonCodePrefix = "c" // 通用模块响应码前缀
)

// 起始编号
const (
	BaseInitCode BaseCode = 0
)

const (
	/* --------------------- 通用模块 --------------------- */

	CodeOK                   BaseCode = iota // 通用：请求成功
	CodeErr                                  // 通用：请求失败
	CodeInternalServerError                  // 通用：系统内部错误
	CodeInvalidRequestParams                 // 通用：请求参数不合法
	CodeUnauthorizedAccess                   // 通用：未授权访问（未登录或 Token 失效）
	CodeAccessDenied                         // 通用：无权限操作

	/* --------------------- 用户模块 --------------------- */

	CodeUserLoginSuccess            // 用户登录成功
	CodeUserLoginFailed             // 用户登录失败（账号或密码错误）
	CodeUserLoginParamError         // 用户参数错误
	CodeUserLoginPasswordMissing    // 用户登录密码未提供（必须填写密码）
	CodeUserLoginIdentifierMissing  // 用户登录用户名或邮箱未提供（必须填写其一）
	CodeUserNotFound                // 用户不存在
	CodeUserEmailNotFound           // 邮箱未注册
	CodeUserDisabled                // 用户账号被禁用
	CodeUserLocked                  // 用户账号被锁定（如连续登录失败）
	CodeUserInputInvalid            // 用户输入参数无效（通用验证失败）
	CodeUserNotVerified             // 用户未完成邮箱或手机验证
	CodeUserLoginAttemptExceeded    // 登录尝试次数过多，触发保护
	CodeUserSessionExpired          // 用户会话过期（AccessToken 过期）
	CodeUserPasswordHashError       // 密码哈希比对失败
	CodeUserRegisterSuccess         // 用户注册成功
	CodeUserUsernameExists          // 用户名已存在
	CodeUserEmailExists             // 邮箱已注册
	CodeUserRegisterFailed          // 用户注册失败（数据库异常等）
	CodeUserVerificationCodeInvalid // 验证码无效
	CodeUserVerificationCodeExpired // 验证码过期
	CodeUserProfileUpdateSuccess    // 用户信息更新成功
	CodeUserProfileUpdateFailed     // 用户信息更新失败
	CodeUserProfileIncomplete       // 用户资料未完善
	CodeUserNicknameTooLong         // 昵称长度超出限制
	CodeUserEmailFormatInvalid      // 邮箱格式非法
	CodeUserPhoneFormatInvalid      // 手机号格式非法

	/* --------------------- 用户 password --------------------- */

	CodeUserPasswordIncorrect          // 密码错误（用户输入不对）
	CodeUserPasswordVerifyError        // 密码校验错误（哈希异常等）
	CodeUserPasswordVerificationFailed // 密码校验失败（统一处理入口）
	CodeUserPasswordHashCompareFailed  // 密码哈希比对失败
	CodeUserPasswordTooWeak            // 密码过于简单（如全数字）
	CodeUserPasswordFormatInvalid      // 密码格式非法（长度/字符校验）
	CodeUserOldPasswordIncorrect       // 原密码错误（用于修改密码时验证）

	/* --------------------- 用户 Token --------------------- */

	CodeTokenGenerationFailed        // Token 生成失败（JWT 过程异常）
	CodeRefreshTokenGenerationFailed // Refresh Token 生成失败
	CodeAccessTokenExpired           // AccessToken 已过期
	CodeRefreshTokenExpired          // RefreshToken 已过期
	CodeTokenInvalidFormat           // Token 无效（格式非法或无法解析）
	CodeRefreshTokenMismatch         // Refresh Token 与当前用户不一致
	CodeAccessTokenStorageFailed     // Access Token 存储失败（例如 Redis 写入失败）
	CodeRefreshTokenStorageFailed    // Refresh Token 存储失败

	/* --------------------- 用户 Redis 其他 --------------------- */

	CodeUserRedisDeviceBindingFailed // 用户设备绑定失败（记录设备信息时出错）

	/* --------------------- Redis 模块 --------------------- */

	CodeRedisConnectionFailed // Redis 连接失败
	CodeRedisReadFailed       // 从 Redis 获取数据失败
	CodeRedisWriteFailed      // 写入 Redis 数据失败
	CodeRedisDeleteFailed     // 删除 Redis 数据失败

	/* --------------------- 数据库模块 --------------------- */

	CodeDatabaseConnectFailed // 数据库连接失败
	CodeDatabaseQueryFailed   // 数据库查询失败
	CodeDatabaseInsertFailed  // 数据插入失败
	CodeDatabaseUpdateFailed  // 数据更新失败
	CodeDatabaseDeleteFailed  // 数据删除失败
)

// 成功类（200 开头，或用户模块成功）
//const (
//	CodeOK                       = BaseInitCode + iota // 通用：请求成功
//	CodeUserLoginSuccess                               // 用户登录成功
//	CodeUserRegisterSuccess                            // 用户注册成功
//	CodeUserProfileUpdateSuccess                       // 用户信息更新成功
//
//	/* --------------------- FailInitCode --------------------- */
//
//	FailInitCode // 失败初始化码
//)
//
//// 用户输入/业务失败类（4xx）
//const (
//	CodeUserLoginFailed         = FailInitCode + 1 + iota // 用户登录失败（账号或密码错误）
//	CodeUserNotFound                                      // 用户不存在
//	CodeUserEmailNotFound                                 // 邮箱未注册
//	CodeUserUsernameExists                                // 用户名已存在
//	CodeUserEmailExists                                   // 邮箱已注册
//	CodeUserRegisterFailed                                // 用户注册失败
//	CodeUserProfileUpdateFailed                           // 用户信息更新失败
//	CodeUserProfileIncomplete                             // 用户资料未完善
//	CodeUserNicknameTooLong                               // 昵称过长
//	CodeUserEmailFormatInvalid                            // 邮箱格式非法
//	CodeUserPhoneFormatInvalid                            // 手机号格式非法
//
//	CodeUserPasswordIncorrect     // 密码错误
//	CodeUserPasswordTooWeak       // 密码过于简单
//	CodeUserPasswordFormatInvalid // 密码格式非法
//	CodeUserOldPasswordIncorrect  // 原密码错误
//
//	CodeUserVerificationCodeInvalid // 验证码无效
//	CodeUserVerificationCodeExpired // 验证码过期
//
//	CodeTokenInvalidFormat   // Token 格式非法
//	CodeRefreshTokenMismatch // Refresh Token 与用户不一致
//
//	/* --------------------- ErrInitCode --------------------- */
//
//	ErrInitCode // 错误初始化码
//)
//
//// 系统内部错误（5xx）
//const (
//	CodeInternalServerError           = ErrInitCode + 1 + iota // 系统内部错误
//	CodeUserPasswordVerifyError                                // 密码校验异常
//	CodeUserPasswordHashCompareFailed                          // 密码哈希比对失败
//
//	CodeTokenGenerationFailed        // Token 生成失败
//	CodeRefreshTokenGenerationFailed // Refresh Token 生成失败
//	CodeAccessTokenStorageFailed     // Access Token 存储失败
//	CodeRefreshTokenStorageFailed    // Refresh Token 存储失败
//
//	CodeUserRedisDeviceBindingFailed // 用户设备绑定失败
//
//	CodeRedisConnectionFailed // Redis 连接失败
//	CodeRedisReadFailed       // Redis 获取失败
//	CodeRedisWriteFailed      // Redis 写入失败
//	CodeRedisDeleteFailed     // Redis 删除失败
//
//	CodeDatabaseConnectFailed // 数据库连接失败
//	CodeDatabaseQueryFailed   // 查询失败
//	CodeDatabaseInsertFailed  // 插入失败
//	CodeDatabaseUpdateFailed  // 更新失败
//	CodeDatabaseDeleteFailed  // 删除失败
//
//	/* --------------------- OtherInitCode --------------------- */
//
//	OtherInitCode // 其他初始化码
//)
