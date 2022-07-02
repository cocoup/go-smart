package sqlx

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type SqlConf struct {
	IP           string `yaml:"ip"`             // 服务器地址
	Port         string `yaml:"port"`           // 端口
	DB           string `yaml:"db"`             // 数据库名
	Name         string `yaml:"username"`       // 数据库用户名
	Password     string `yaml:"password"`       // 数据库密码
	Option       string `yaml:"Option"`         // 高级配置
	MaxIdleConns int    `yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `yaml:"log-mode"`       // 是否开启Gorm全局日志
}

func (m *SqlConf) DNS() string {
	return m.Name + ":" + m.Password + "@tcp(" + m.IP + ":" + m.Port + ")/" + m.DB + "?" + m.Option
}

func gormConfig(conf SqlConf) *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}
	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})

	switch conf.LogMode {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
