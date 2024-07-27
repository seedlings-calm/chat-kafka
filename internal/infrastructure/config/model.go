package config

type BaseViper struct {
	Redis    []Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Kafka    Kafka      `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
	Postgres []Postgres `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Name     string
}

type Redis struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"` // 服务器地址
	Port int    `json:"port" yaml:"port" mapstructure:"port"` //端口
	Pwd  string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`    // 密码
	Db   int    `mapstructure:"db" json:"db" yaml:"db"`       // 单实例模式下redis的哪个数据库
	Name string `mapstructure:"name" json:"name" yaml:"name"` //
}

type Kafka struct {
	Private   KafkaConf //私聊
	Group     KafkaConf //群聊
	Broadcast KafkaConf //广播
}

type KafkaConf struct {
	Addr  []string `mapstructure:"addr" json:"addr" yaml:"addr"`    //服务地址
	Group string   `mapstructure:"group" json:"group" yaml:"group"` //组
	Topic string   `mapstructure:"topic" json:"topic" yaml:"topic"` //主题
}
