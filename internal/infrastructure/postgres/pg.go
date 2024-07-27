package postgres

import (
	"fmt"
	"sync"

	"github.com/seedlings-calm/chat-kafka/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	sets *PgSqlSets
)

type PgSqlSets struct {
	pgsql map[string]*gorm.DB
	l     sync.RWMutex
}

func (r *PgSqlSets) Db(key ...string) *gorm.DB {
	r.l.RLock()
	defer r.l.RUnlock()
	var name string
	if len(key) <= 0 {
		name = "default"
	} else {
		name = key[0]
	}
	if db, ok := r.pgsql[name]; ok {
		return db
	}
	return nil
}

func NewDb() *PgSqlSets {
	once.Do(func() {
		conf := config.GetGlobalConf()
		dbClients := map[string]*gorm.DB{}
		for _, m := range conf.Postgres {
			dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", m.Host, m.User, m.Password, m.Database, m.Port)
			open, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				panic("pgsql数据库初始化失败:" + err.Error())
			}
			db, err := open.DB()
			if err != nil {
				panic(err)
			}
			db.SetMaxOpenConns(5)
			db.SetMaxIdleConns(10)
			dbClients[m.Name] = open
		}
		sets = &PgSqlSets{
			pgsql: dbClients,
		}

	})
	return sets
}
