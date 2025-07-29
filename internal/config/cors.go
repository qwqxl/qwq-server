package config

/*
   cors:
       enabled: true
       allow_origins: [ "*" ]
       allow_methods: [ "GET", "POST", "PUT", "DELETE", "OPTIONS" ]
       allow_headers: [ "Origin", "Content-Type", "Authorization", "Accept" ]
       expose_headers: [ "Content-Length" ]
       max_age: 86400
       allow_credentials: true
*/

// Cors 跨域配置
type Cors struct {
	Enabled             bool     `mapstructure:"enabled" yaml:"enabled" qwq-default:"true"`
	AllowOrigins        []string `mapstructure:"allow_origins" yaml:"allow_origins" qwq-default:"*"`
	AllowOriginPatterns []string `mapstructure:"allow_origin_patterns" yaml:"allow_origin_patterns"`
	AllowMethods        []string `mapstructure:"allow_methods" yaml:"allow_methods" qwq-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowHeaders        []string `mapstructure:"allow_headers" yaml:"allow_headers" qwq-default:"Origin,Content-Type,Authorization,Accept"`
	ExposeHeaders       []string `mapstructure:"expose_headers" yaml:"expose_headers" qwq-default:"Content-Length"`
	MaxAge              int      `mapstructure:"max_age" yaml:"max_age" qwq-default:"86400"`
	AllowCredentials    bool     `mapstructure:"allow_credentials" yaml:"allow_credentials" qwq-default:"true"`
	Debug               bool     `mapstructure:"debug" yaml:"debug" qwq-default:"true"`
}
