package config

type Security struct {
	CSRF CSRF `mapstructure:"csrf" yaml:"csrf"`
	XSS  XSS  `mapstructure:"xss" yaml:"xss"`
}

type CSRF struct {
	Enabled        bool   `mapstructure:"enabled" yaml:"enabled" qwq-default:"true"`                   // 是否启用CSRF
	Secret         string `mapstructure:"secret" yaml:"secret" qwq-default:"your-32-byte-secret-key"`  // 实际使用中应替换为随机密钥
	TokenLength    int    `mapstructure:"token_length" yaml:"token_length" qwq-default:"32"`           // CSRF令牌长度
	CookieName     string `mapstructure:"cookie_name" yaml:"cookie_name" qwq-default:"_csrf_token"`    // Cookie名称
	CookieSecure   bool   `mapstructure:"cookie_secure" yaml:"cookie_secure" qwq-default:"true"`       // Cookie是否启用安全传输
	CookieHTTPOnly bool   `mapstructure:"cookie_http_only" yaml:"cookie_http_only" qwq-default:"true"` // Cookie是否仅允许HTTP访问
	CookieMaxAge   int    `mapstructure:"cookie_max_age" yaml:"cookie_max_age" qwq-default:"86400"`    // Cookie最大有效期
}

type XSS struct {
	Enabled               bool   `mapstructure:"enabled" yaml:"enabled" qwq-default:"true"`
	XSSProtection         string `mapstructure:"x_xss_protection" yaml:"x_xss_protection" qwq-default:"1; mode=block"`
	ContentTypeNosniff    string `mapstructure:"content_type_nosniff" yaml:"content_type_nosniff" qwq-default:"nosniff"`
	XFrameOptions         string `mapstructure:"x_frame_options" yaml:"x_frame_options" qwq-default:"DENY"`
	HSTSMaxAge            int    `mapstructure:"hsts_max_age" yaml:"hsts_max_age" qwq-default:"31536000"`
	HSTSIncludeSubdomains bool   `mapstructure:"hsts_include_subdomains" yaml:"hsts_include_subdomains" qwq-default:"true"`
	ContentSecurityPolicy string `mapstructure:"content_security_policy" yaml:"content_security_policy" qwq-default:"default-src 'self'"`
}
