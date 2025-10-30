package errno

const (
	SuccessCode = 10000
	SuccessMsg  = "ok"
)

// 200xx: 参数错误，Param 打头
const (
	ParamVerifyErrorCode   = 20000 + iota // 参数校验失败
	ParamMissingErrorCode                 // 参数缺失
	ParamMissingHeaderCode                // 请求头缺失
	ParamInvalidHeaderCode                // 请求头无效
)

// 300xx: 鉴权错误，Auth 打头
const (
	AuthInvalidCode             = 30000 + iota // 鉴权失败
	AuthAccessExpiredCode                      // 访问令牌过期
	AuthRefreshExpiredCode                     // 刷新令牌过期
	AuthNoTokenCode                            // 没有 token
	AuthNoOperatePermissionCode                // 没有操作权限
	AuthMissingTokenCode                       // 缺少 token
	IllegalOperatorCode                        // 不合格的操作(比如传入 payment status时传入了一个不存在的 status)
)

// 500xx: 内部错误，Internal 打头
// 服务级别的错误, 发生的时候说明我们程序自身出了问题
// 比如数据库断联, 编码错误等. 需要我们人为的去维护
const (
	InternalServiceErrorCode  = 50000 + iota // 内部服务错误
	InternalDatabaseErrorCode                // 数据库错误
	InternalRedisErrorCode                   // Redis错误
	InternalNetworkErrorCode                 // 网络错误
	OSOperateErrorCode
	IOOperateErrorCode
	InsufficientStockErrorCode
	InternalRPCErrorCode
	InternalRocketmqErrorCode
)

const (
	UpYunFileErrorCode = 60000 + iota
	RedisKeyNotExist
	RepeatedOperation
)

// User
const (
	ServiceUserExist = 1000 + iota
	ServiceUserNotExist
	ServiceInvalidUsername
	ServiceInvalidPassword
	ServiceInvalidEmail
	ServiceInvalidCode

	ErrRecordNotFound
	UserLogOut
	UserAlreadyLogin
)
