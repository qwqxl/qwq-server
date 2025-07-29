package authx

// HookFunc 定义了生命周期钩子函数的类型
type HookFunc func(userID, platform, device string) error

// LifecycleHooks 包含了所有生命周期钩子
type LifecycleHooks struct {
	OnLogin   HookFunc // 登录时触发
	OnLogout  HookFunc // 登出时触发
	OnRefresh HookFunc // 刷新令牌时触发
	OnKick    HookFunc // 用户被踢出时触发
}