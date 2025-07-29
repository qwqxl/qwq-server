# RBAC 

## 使用示例

初始化 RBAC

```go
package main

import (
	"fmt"
	"yourproject/pkg/rbac"
)

func main() {
	// 创建新的 RBAC 实例
	rbacSystem := rbac.NewRBAC()
	
	// 添加权限
	rbacSystem.AddPermission("create_post")
	rbacSystem.AddPermission("edit_post")
	rbacSystem.AddPermission("delete_post")
	rbacSystem.AddPermission("view_post")
	
	// 添加角色
	editorPerms := []rbac.Permission{"create_post", "edit_post", "view_post"}
	if err := rbacSystem.AddRole("editor", editorPerms); err != nil {
		fmt.Println("Error adding editor role:", err)
	}
	
	adminPerms := []rbac.Permission{"create_post", "edit_post", "delete_post", "view_post"}
	if err := rbacSystem.AddRole("admin", adminPerms); err != nil {
		fmt.Println("Error adding admin role:", err)
	}
	
	// 添加用户
	if err := rbacSystem.AddUser("user1"); err != nil {
		fmt.Println("Error adding user1:", err)
	}
	if err := rbacSystem.AddUser("user2"); err != nil {
		fmt.Println("Error adding user2:", err)
	}
	
	// 为用户分配角色
	if err := rbacSystem.AssignRole("user1", "editor"); err != nil {
		fmt.Println("Error assigning role to user1:", err)
	}
	if err := rbacSystem.AssignRole("user2", "admin"); err != nil {
		fmt.Println("Error assigning role to user2:", err)
	}
	
	// 检查权限
	checkPermission(rbacSystem, "user1", "delete_post") // 应该返回 false
	checkPermission(rbacSystem, "user2", "delete_post") // 应该返回 true
	
	// 为编辑器添加删除权限
	if err := rbacSystem.AddPermissionToRole("editor", "delete_post"); err != nil {
		fmt.Println("Error adding permission to editor:", err)
	}
	
	checkPermission(rbacSystem, "user1", "delete_post") // 现在应该返回 true
	
	// 移除用户的角色
	if err := rbacSystem.RemoveRoleFromUser("user1", "editor"); err != nil {
		fmt.Println("Error removing role from user1:", err)
	}
	
	checkPermission(rbacSystem, "user1", "delete_post") // 应该返回 false
	
	// 获取用户角色
	if roles, err := rbacSystem.GetUserRoles("user2"); err == nil {
		fmt.Println("User2 roles:", roles) // 应该输出 [admin]
	}
	
	// 获取角色权限
	if perms, err := rbacSystem.GetRolePermissions("admin"); err == nil {
		fmt.Println("Admin permissions:", perms)
	}
}

func checkPermission(rbacSystem *rbac.RBAC, userID rbac.UserID, perm rbac.Permission) {
	hasPerm, err := rbacSystem.HasPermission(userID, perm)
	if err != nil {
		fmt.Printf("Error checking permission for %s: %v\n", userID, err)
		return
	}
	
	fmt.Printf("User %s has permission '%s': %t\n", userID, perm, hasPerm)
}
```

## 持久化集成

```go
package main

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
	"yourproject/pkg/rbac"
)

var rbacSystem *rbac.RBAC

func main() {
	// 初始化RBAC系统
	rbacSystem = rbac.NewRBAC()
	
	// 添加默认权限和角色
	initializeRBAC()
	
	// 创建Gin路由
	router := gin.Default()
	
	// 添加权限检查中间件
	router.Use(authMiddleware)
	
	// 添加RBAC管理端点
	rbacGroup := router.Group("/rbac")
	{
		rbacGroup.POST("/permissions", addPermissionHandler)
		rbacGroup.POST("/roles", addRoleHandler)
		rbacGroup.POST("/users", addUserHandler)
		rbacGroup.POST("/assign", assignRoleHandler)
	}
	
	// 添加需要权限保护的端点
	protectedGroup := router.Group("/api")
	{
		protectedGroup.POST("/posts", checkPermission("create_post"), createPostHandler)
		protectedGroup.PUT("/posts/:id", checkPermission("edit_post"), updatePostHandler)
		protectedGroup.DELETE("/posts/:id", checkPermission("delete_post"), deletePostHandler)
	}
	
	router.Run(":8080")
}

func initializeRBAC() {
	// 添加权限
	rbacSystem.AddPermission("create_post")
	rbacSystem.AddPermission("edit_post")
	rbacSystem.AddPermission("delete_post")
	rbacSystem.AddPermission("view_post")
	
	// 添加角色
	editorPerms := []rbac.Permission{"create_post", "edit_post", "view_post"}
	rbacSystem.AddRole("editor", editorPerms)
	
	adminPerms := []rbac.Permission{"create_post", "edit_post", "delete_post", "view_post"}
	rbacSystem.AddRole("admin", adminPerms)
}

// 权限检查中间件
func checkPermission(perm rbac.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID（假设已通过认证）
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		
		// 检查权限
		hasPerm, err := rbacSystem.HasPermission(rbac.UserID(userID.(string)), perm)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			return
		}
		
		if !hasPerm {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		
		c.Next()
	}
}

// 示例：添加权限端点
func addPermissionHandler(c *gin.Context) {
	var request struct {
		Permission string `json:"permission"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	
	rbacSystem.AddPermission(rbac.Permission(request.Permission))
	c.JSON(http.StatusOK, gin.H{"status": "Permission added"})
}
```