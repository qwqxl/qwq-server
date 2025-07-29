package addresslock

import (
	"sync"
	"sync/atomic"
)

// 地址锁条目，包含读写锁和引用计数
type lockEntry struct {
	mu       sync.RWMutex
	refCount int32 // 当前使用该锁的协程数
	removed  int32 // 标记是否已移除 (0: 未移除, 1: 已移除)
}

// 多地址锁管理器
type AddressLock struct {
	locks sync.Map // 存储地址到lockEntry的映射 (key: interface{}, value: *lockEntry)
}

// 获取指定地址的锁条目（内部方法）
func (al *AddressLock) getOrCreateEntry(key interface{}) *lockEntry {
	for {
		// 尝试获取现有条目
		if actual, loaded := al.locks.Load(key); loaded {
			entry := actual.(*lockEntry)

			// 原子增加引用计数
			atomic.AddInt32(&entry.refCount, 1)

			// 检查条目是否有效
			if atomic.LoadInt32(&entry.removed) == 0 {
				return entry
			}

			// 条目已失效，回退引用计数并重试
			atomic.AddInt32(&entry.refCount, -1)
		}

		// 创建新条目
		newEntry := &lockEntry{refCount: 1}
		actual, loaded := al.locks.LoadOrStore(key, newEntry)

		if !loaded {
			return newEntry // 成功创建新条目
		}

		// 其他协程已创建，使用现有条目
		entry := actual.(*lockEntry)
		atomic.AddInt32(&entry.refCount, 1)
		if atomic.LoadInt32(&entry.removed) == 0 {
			return entry
		}

		// 条目已失效，回退引用计数并重试
		atomic.AddInt32(&entry.refCount, -1)
	}
}

// 释放锁条目引用（内部方法）
func (al *AddressLock) releaseEntry(key interface{}, entry *lockEntry) {
	// 减少引用计数
	if atomic.AddInt32(&entry.refCount, -1) == 0 {
		// 标记为已移除
		if atomic.CompareAndSwapInt32(&entry.removed, 0, 1) {
			// 从映射中安全移除
			al.locks.CompareAndDelete(key, entry)
		}
	}
}

// 获取指定地址的写锁
func (al *AddressLock) Lock(key interface{}) {
	entry := al.getOrCreateEntry(key)
	entry.mu.Lock()
}

// 释放写锁
func (al *AddressLock) Unlock(key interface{}) {
	if actual, ok := al.locks.Load(key); ok {
		entry := actual.(*lockEntry)
		entry.mu.Unlock()
		al.releaseEntry(key, entry)
	} else {
		panic("unlock of unlocked mutex")
	}
}

// 获取指定地址的读锁
func (al *AddressLock) RLock(key interface{}) {
	entry := al.getOrCreateEntry(key)
	entry.mu.RLock()
}

// 释放读锁
func (al *AddressLock) RUnlock(key interface{}) {
	if actual, ok := al.locks.Load(key); ok {
		entry := actual.(*lockEntry)
		entry.mu.RUnlock()
		al.releaseEntry(key, entry)
	} else {
		panic("runlock of unlocked mutex")
	}
}

// 使用示例
// var al AddressLock
// al.Lock("account:123")
// // 执行写操作
// al.Unlock("account:123")
//
// al.RLock("account:456")
// // 执行读操作
// al.RUnlock("account:456")
