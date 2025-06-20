package perm

type Permission uint64

/* Permission 基础权限定义 */

const None Permission = 0 // 无权限
const (
	/* User 用户管理权限 */

	UserProfileView    Permission = 1 << iota // 查看用户资料
	UserProfileEditOwn                        // 编辑自己的资料
	UserProfileEditAny                        // 编辑任意用户资料
	UserBan                                   // 封禁/解封用户
	UserWarn                                  // 发送警告

	/* Post 帖子相关权限 */

	PostRead      // 查看帖子（基础权限）
	PostCreate    // 创建新帖子
	PostEditOwn   // 编辑自己的帖子
	PostDeleteOwn // 删除自己的帖子
	PostPin       // 置顶帖子
	PostLock      // 锁定/解锁帖子
	PostDeleteAny // 删除任意帖子

	/* Comment 评论相关权限 */

	CommentCreate    // 发表评论
	CommentEditOwn   // 编辑自己的评论
	CommentDeleteOwn // 删除自己的评论
	CommentDeleteAny // 删除任意评论

	/* Board 版块管理权限 */

	BoardCreate       // 创建新版块
	BoardModify       // 修改版块信息
	BoardDelete       // 删除版块
	BoardManageAccess // 管理版块访问权限

	/* Content 内容管理权限 */

	ContentAudit        // 审核待审内容
	ContentFeature      // 设置精华/推荐内容
	ContentReportView   // 查看举报内容
	ContentReportManage // 处理举报内容

	/* System 系统管理权限 */

	SysConfig  // 修改系统配置
	SysBackup  // 执行系统备份
	SysLogView // 查看系统日志

	/* PrivateMessage 私信权限 */

	PMSend   // 发送私信
	PMRead   // 查看私信
	PMDelete // 删除私信

	/* Attachment 附件权限 */

	AttachmentUpload    // 上传附件
	AttachmentDownload  // 下载附件
	AttachmentDeleteOwn // 删除自己的附件
	AttachmentDeleteAny // 删除任意附件

	/* Other 其他权限 */

	permMax // 权限校验边界
)

type PermissionGroup struct {
	Permission Permission // 权限组
}
