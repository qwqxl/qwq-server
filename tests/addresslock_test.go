package qwqtest

import (
	"fmt"
	"qwqserver/pkg/addresslock"
	"sync"
	"testing"
)

type ALUser struct {
	ID    int
	Name  string
	age   int
	Money int
}

type UsersManages struct {
	al    addresslock.AddressLock
	users map[int]*ALUser // 使用map存储用户，ID作为键
	mu    sync.RWMutex    // 保护整个users映射
}

func NewUsersManages() *UsersManages {
	return &UsersManages{
		users: make(map[int]*ALUser),
	}
}

func (m *UsersManages) Add(user *ALUser) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
}

func (m *UsersManages) GetUser(ID int) (*ALUser, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, exists := m.users[ID]
	return user, exists
}

// Transfer 安全转账操作
func (m *UsersManages) Transfer(fromID int, toID int, amount int) error {
	// 确定锁定顺序，避免死锁
	firstID, secondID := fromID, toID
	if fromID > toID {
		firstID, secondID = toID, fromID
	}

	// 锁定两个账户
	m.al.Lock(firstID)
	defer m.al.Unlock(firstID)

	if firstID != secondID {
		m.al.Lock(secondID)
		defer m.al.Unlock(secondID)
	}

	// 获取用户
	fromUser, exists := m.GetUser(fromID)
	if !exists {
		return fmt.Errorf("sender %d not found", fromID)
	}

	toUser, exists := m.GetUser(toID)
	if !exists {
		return fmt.Errorf("receiver %d not found", toID)
	}

	// 检查余额
	if fromUser.Money < amount {
		return fmt.Errorf("insufficient funds in account %d", fromID)
	}

	// 执行转账
	fromUser.Money -= amount
	toUser.Money += amount

	return nil
}

// GetBalance 获取用户余额
func (m *UsersManages) GetBalance(ID int) (int, error) {
	m.al.RLock(ID)
	defer m.al.RUnlock(ID)

	user, exists := m.GetUser(ID)
	if !exists {
		return 0, fmt.Errorf("user %d not found", ID)
	}
	return user.Money, nil
}

// UpdateUser 更新用户信息
func (m *UsersManages) UpdateUser(user *ALUser) error {
	m.al.Lock(user.ID)
	defer m.al.Unlock(user.ID)

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; !exists {
		return fmt.Errorf("user %d not found", user.ID)
	}

	// 实际项目中这里会有更复杂的更新逻辑
	m.users[user.ID] = user
	return nil
}

// BatchTransfer 批量转账操作
func (m *UsersManages) BatchTransfer(transfers []struct{ From, To, Amount int }) error {
	// 按ID排序转账请求以避免死锁
	sortedTransfers := make([]struct{ From, To, Amount int }, len(transfers))
	copy(sortedTransfers, transfers)

	// 简单排序：确保每个转账中的ID按顺序处理
	for i := range sortedTransfers {
		if sortedTransfers[i].From > sortedTransfers[i].To {
			sortedTransfers[i].From, sortedTransfers[i].To = sortedTransfers[i].To, sortedTransfers[i].From
			sortedTransfers[i].Amount = -sortedTransfers[i].Amount // 反转金额方向
		}
	}

	// 按FromID和ToID排序整个列表
	// 实际项目中应使用更高效的排序算法
	for i := 0; i < len(sortedTransfers); i++ {
		for j := i + 1; j < len(sortedTransfers); j++ {
			if sortedTransfers[i].From > sortedTransfers[j].From ||
				(sortedTransfers[i].From == sortedTransfers[j].From &&
					sortedTransfers[i].To > sortedTransfers[j].To) {
				sortedTransfers[i], sortedTransfers[j] = sortedTransfers[j], sortedTransfers[i]
			}
		}
	}

	// 执行批量转账
	for _, t := range sortedTransfers {
		if err := m.Transfer(t.From, t.To, t.Amount); err != nil {
			return err
		}
	}
	return nil
}

func TestAddressLockMain(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		alUserTest()
	})
}

func alUserTest() {
	manager := NewUsersManages()

	// 添加初始用户
	manager.Add(&ALUser{ID: 1, Name: "Alice", Money: 1000})
	manager.Add(&ALUser{ID: 2, Name: "Bob", Money: 500})
	manager.Add(&ALUser{ID: 3, Name: "Charlie", Money: 200})

	var wg sync.WaitGroup

	// 并发转账
	wg.Add(3)
	go func() {
		defer wg.Done()
		manager.Transfer(1, 2, 100) // Alice 转给 Bob 100
	}()
	go func() {
		defer wg.Done()
		manager.Transfer(2, 3, 50) // Bob 转给 Charlie 50
	}()
	go func() {
		defer wg.Done()
		manager.Transfer(3, 1, 30) // Charlie 转给 Alice 30
	}()

	wg.Wait()

	// 查询余额
	aliceBalance, _ := manager.GetBalance(1)
	bobBalance, _ := manager.GetBalance(2)
	charlieBalance, _ := manager.GetBalance(3)

	fmt.Printf("Final balances:\n")
	fmt.Printf("Alice (%d): $%d\n", 1, aliceBalance)
	fmt.Printf("Bob (%d): $%d\n", 2, bobBalance)
	fmt.Printf("Charlie (%d): $%d\n", 3, charlieBalance)
	fmt.Printf("Total money: $%d\n", aliceBalance+bobBalance+charlieBalance)
}
