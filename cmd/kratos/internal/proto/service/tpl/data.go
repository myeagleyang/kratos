package tpl

var (
	DataTemplate = `package data

import (
	"github.com/IBM/sarama"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/google/wire"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/conf"
	dbPkg "gitlab.wwgame.com/chaoshe/blind_box/pkg/gormdb"
	"gitlab.wwgame.com/wwgame/kit/util"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, New{{ .Service }}Repo)

// Data .
type Data struct {
	// TODO wrapped database client
	db    *gorm.DB
	cache *redis.Client
	rs    *redsync.Redsync
	kap   sarama.AsyncProducer
	kcg   sarama.ConsumerGroup
}

// NewData .
func NewData(c *conf.Data, lg log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(lg).Info("closing the data resources")
	}
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{Logger: logger.Default.LogMode(logger.LogLevel(c.Database.LogLevel))})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Duration(100) * time.Second)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(300)
	// 迁移表
	err = dbPkg.MigrateTable(db, []dbPkg.Model{
		// 添加迁移模型
	})
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Network:      c.Redis.Network,
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           util.GetInt(c.Redis.Database),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		PoolSize:     util.GetInt(c.Redis.PoolSize),
	})
	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	producer := NewKafkaProducer(c)
	// kafka client
	consumer := NewKafkaConsumerGroup(c)
	return &Data{db: db, cache: client, rs: rs, kap: producer, kcg: consumer}, cleanup, nil
}

func (d *Data) GetDB() *gorm.DB {
	return d.db
}

func (d *Data) NewRedisMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return d.rs.NewMutex(name, options...)
}

`
	DataImplTemplate = `package data

import (
	//"context"

	//pb "gitlab.wwgame.com/chaoshe/blind_box/api/{{ .ServiceLower }}/{{ .Mode }}/v1"
	"gitlab.wwgame.com/chaoshe/blind_box/app/{{ .ServiceLower }}/{{ .Mode }}/internal/biz"
	"gitlab.wwgame.com/wwgame/kratos/v2/log"
)

type {{ .ServiceLower }}Repo struct {
	data *Data
	log  *log.Helper
}

func New{{ .Service }}Repo(data *Data, logger log.Logger) biz.{{ .Service }}Repo {
	return &{{ .ServiceLower }}Repo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *{{ .ServiceLower }}Repo) GetRedisMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return r.data.NewRedisMutex(name, options...)
}

func (r *{{ .ServiceLower }}Repo) StartTx() *gorm.DB {
	return r.data.db.Begin()
}

func (r *{{ .ServiceLower }}Repo) GetDB() *gorm.DB {
	return r.data.db
}
`
)
