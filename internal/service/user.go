package service

import (
	"context"
	"qwqserver/internal/model"
	"qwqserver/internal/repository"
)

// UserDelete 删除用户
func UserDelete(uid uint) (res *model.Result) {
	// 初始化 返回结果
	res = &model.Result{}
	user := &model.User{}

	if uid == 0 {
		res.Code = 400
		res.Message = "用户ID不能为: 0"
		return
	}

	userRepo, err := repository.NewUserRepository()
	if err != nil {
		res.Code = 500
		res.Message = "获取数据库连接失败: " + err.Error()
		return
	}

	// 检查用户是否存在
	if user, err = userRepo.FindByID(context.Background(), uid); err != nil || user == nil {
		res.Code = 404
		res.Message = "用户不存在"
		return
	}

	if err = userRepo.Delete(context.Background(), user.ID); err != nil {
		res.Code = 500
		res.Message = "删除用户失败: " + err.Error()
		return
	}

	res.Code = 200
	res.Message = "删除用户成功"
	res.Data = map[string]any{
		"uid":      user.ID,
		"username": user.Username,
		"email":    user.Email,
	}
	return
}
