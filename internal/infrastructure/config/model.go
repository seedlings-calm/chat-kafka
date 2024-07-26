package config

type BaseViper struct {
	Redis []Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type Redis struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"` // 服务器地址
	Port int    `json:"port" yaml:"port" mapstructure:"port"` //端口
	Pwd  string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`    // 密码
	Db   int    `mapstructure:"db" json:"db" yaml:"db"`       // 单实例模式下redis的哪个数据库
	Name string `mapstructure:"name" json:"name" yaml:"name"` //
}
