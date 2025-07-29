package perm

import (
	"fmt"
	"testing"
)

func TestPermRBAC(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		rbac := NewRBAC()

		rbac.AddPermission("read_article")
		rbac.AddPermission("write_article")

		rbac.AddRole("reader", []Permission{"read_article"})
		rbac.AddRole("writer", []Permission{"read_article", "write_article"})

		rbac.AddUser("user1")
		rbac.AssignRole("user1", "reader")

		rbac.AddUser("user2")
		rbac.AssignRole("user2", "writer")

		fmt.Println(rbac.HasPermission("user1", "read_article"))  // true
		fmt.Println(rbac.HasPermission("user1", "write_article")) // false
		fmt.Println(rbac.HasPermission("user2", "write_article")) // true
	})

}
