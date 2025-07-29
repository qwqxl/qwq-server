package plugin

import (
	"qwqserver/pkg/util/singleton"
	"strings"
	"sync"
)

// RegisterRequest is the request for registering a plugin
type RegisterRequest struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type RegisterMap struct {
	Plugins map[string]string
	// 加锁
	mu sync.RWMutex
}

func (p *RegisterMap) Map(name string, addrs ...string) string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(addrs) == 0 {
		//p := NewPluginsMap()
		return p.Plugins[name]
	}
	addr := strings.Join(addrs, "")
	p.Plugins[name] = addr
	return addr
}

// PluginsMapData 单例
var PluginsMapData singleton.Singleton[RegisterMap]

// NewPluginsMap Register registers a plugin
func NewPluginsMap() *RegisterMap {
	m, err := PluginsMapData.Get(func() (*RegisterMap, error) {
		return &RegisterMap{
			Plugins: make(map[string]string), // ✅ 这里必须初始化
		}, nil
	})
	if err != nil {
		panic(err)
	}
	return m
}

func RegisterPlugin(req *RegisterRequest) *RegisterResponse {
	res := &RegisterResponse{}
	if req.Name == "" || req.Addr == "" {
		res.Error = "Invalid request"
		return res
	}
	if CheckPluginHealth(req.Name, req.Addr) {
		res.Error = "Plugin already registered"
		return res
	}
	RegisterToMainService(req.Name, req.Addr)
	res.Success = true

	return res
}
