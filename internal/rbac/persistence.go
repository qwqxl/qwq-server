package rbac

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// 可序列化的RBAC状态
type RBACState struct {
	Users       map[UserID]*User `json:"users"`
	Roles       map[RoleID]*Role `json:"roles"`
	Permissions []Permission     `json:"permissions"`
}

// 保存RBAC状态到文件
func (r *RBAC) SaveToFile(filename string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	state := RBACState{
		Users:       r.users,
		Roles:       r.roles,
		Permissions: make([]Permission, 0, len(r.permissions)),
	}

	for p := range r.permissions {
		state.Permissions = append(state.Permissions, p)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// 从文件加载RBAC状态
func LoadFromFile(filename string) (*RBAC, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var state RBACState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	rbac := NewRBAC()

	// 重建权限集合
	for _, p := range state.Permissions {
		rbac.permissions[p] = struct{}{}
	}

	// 重建角色
	rbac.roles = state.Roles

	// 重建用户
	rbac.users = state.Users

	return rbac, nil
}

// 定期自动保存（可选）
func (r *RBAC) StartAutoSave(filename string, interval int) chan struct{} {
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Second):
				if err := r.SaveToFile(filename); err != nil {
					fmt.Printf("Auto-save failed: %v\n", err)
				}
			case <-stop:
				if err := r.SaveToFile(filename); err != nil {
					fmt.Printf("Final save failed: %v\n", err)
				}
				return
			}
		}
	}()

	return stop
}
