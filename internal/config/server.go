package config

type Server struct {
	Mode     string `yaml:"mode" env:"SERVER_MODE" env-default:"debug" qwq-default:"debug"`
	Cors     `mapstructure:"cors" yaml:"cors"`
	Security `yaml:"security"`
	Web      []*Web `yaml:"web"`
}

type Web struct {
	ListenAddr   string         `yaml:"listen_addr" qwq-default:""` // 监听地址 (e.g. ":80", ":443")
	TLS          *TLSConfig     `yaml:"tls,omitempty"`              // TLS 配置
	VirtualHosts []*VirtualHost `yaml:"virtual_hosts"`              // 虚拟主机配置
}

// 添加负载均衡和路由规则配置
type VirtualHost struct {
	Hostname string `yaml:"hostname" qwq-default:"proxy.iqwq.com"`               // 域名
	Proxy    string `yaml:"proxy,omitempty" qwq-default:"http://localhost:8080"` // 向后兼容：单个代理
	RootDir  string `yaml:"root_dir,omitempty" qwq-default:""`                   // 向后兼容：静态资源根目录

	Backends []string `yaml:"backends,omitempty" qwq-default:"[]"`           // 负载均衡后端地址列表
	LBPolicy string   `yaml:"lb_policy,omitempty" qwq-default:"round_robin"` // 负载策略：round_robin / random

	Routes []Route `yaml:"routes,omitempty" qwq-default:"[]"` // 路由规则

	HealthPath  string `yaml:"health_path,omitempty" qwq-default:"/health"` // 健康检查路径
	HealthCheck int    `yaml:"health_interval,omitempty" qwq-default:"15"`  // 健康检查周期（秒）
}

// 路由规则结构
type Route struct {
	Path    string `yaml:"path" qwq-default:"/"`              // 路径前缀（如 /api/）
	Proxy   string `yaml:"proxy,omitempty" qwq-default:""`    // 路由反代目标
	RootDir string `yaml:"root_dir,omitempty" qwq-default:""` // 路由静态目录
}

// TLS 配置结构
type TLSConfig struct {
	Enabled       bool   `yaml:"enabled" qwq-default:"false"`         // 是否启用 TLS
	CertFile      string `yaml:"cert_file" qwq-default:""`            // 证书文件路径
	KeyFile       string `yaml:"key_file" qwq-default:""`             // 私钥文件路径
	MinVersion    string `yaml:"min_version" qwq-default:"TLS12"`     // 最小版本 TLS12 / TLS13
	StrictCiphers bool   `yaml:"strict_ciphers" qwq-default:"true"`   // 是否启用严格加密套件
	HSTSMaxAge    int    `yaml:"hsts_max_age" qwq-default:"63072000"` // HSTS 生效时间（秒）
}

//// 添加负载均衡和路由规则配置
//type VirtualHost struct {
//	Hostname string `yaml:"hostname" qwq-default:"www.iqwq.com"`
//	Proxy    string `yaml:"proxy,omitempty" qwq-default:"http://localhost:8080"` // 向后兼容
//	RootDir  string `yaml:"root_dir,omitempty" qwq-default:""`
//
//	// 新增高级功能
//	Backends []string `yaml:"backends,omitempty"` // 负载均衡后端列表
//	LBPolicy string   `yaml:"lb_policy,omitempty" qwq-default:"round_robin"`
//
//	Routes []Route `yaml:"routes,omitempty"` // 路径路由规则
//
//	HealthPath  string `yaml:"health_path,omitempty" qwq-default:"/health"`
//	HealthCheck int    `yaml:"health_interval,omitempty" qwq-default:"15"`
//}
//
//// 新增路由规则结构
//type Route struct {
//	Path    string `yaml:"path"`               // 路径前缀 (e.g. "/api/")
//	Proxy   string `yaml:"proxy,omitempty"`    // 代理目标
//	RootDir string `yaml:"root_dir,omitempty"` // 静态文件目录
//}
//
//// 在TLS配置中添加高级安全选项
//type TLSConfig struct {
//	Enabled       bool   `yaml:"enabled" qwq-default:"false"`
//	CertFile      string `yaml:"cert_file" qwq-default:""`
//	KeyFile       string `yaml:"key_file" qwq-default:""`
//	MinVersion    string `yaml:"min_version" qwq-default:"TLS12"` // TLS12/TLS13
//	StrictCiphers bool   `yaml:"strict_ciphers" qwq-default:"true"`
//	HSTSMaxAge    int    `yaml:"hsts_max_age" qwq-default:"63072000"` // HSTS有效期
//}
