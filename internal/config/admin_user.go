package config

type AdminUser struct {
	Username string `yaml:"username" env:"ADMIN_USERNAME" env-default:"admin" qwq-default:"admin"`
	Password string `yaml:"password" env:"ADMIN_PASSWORD" env-default:"admin" qwq-default:"admin"`
	Nickname string `yaml:"nickname" env:"ADMIN_NICKNAME" env-default:"管理员" qwq-default:"管理员"`
	Email    string `yaml:"email" env:"ADMIN_EMAIL" env-default:"admin@example.com" qwq-default:"admin@example.com"`
}
