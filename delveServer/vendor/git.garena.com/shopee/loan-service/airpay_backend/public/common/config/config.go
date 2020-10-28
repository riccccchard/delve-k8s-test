package config

import (
	"errors"
	"gopkg.in/ini.v1"
	"strings"
	"sync"
	"time"
)

//框架列出了基础的配置部分,后续根据需要的时候可以继续添加

//BaseDbConfig 数据库基本配置
type BaseDbConfig struct {
	DbHost        string `ini:"db_host"`
	DbUser        string `ini:"db_user"`
	DbPassword    string `ini:"db_password"`
	DbDatabase    string `ini:"db_database"`
	DbPoolSize    int    `ini:"db_pool_size"`
	DbMaxIdleSize int    `ini:"db_max_idle_size"`
	DbShowSql     bool   `ini:"db_show_sql"`
}

//DbConfig 数据库主从配置
type DbConfig struct {
	Master BaseDbConfig
	Slave  []BaseDbConfig
}

//RedisConfig  go-redis需要的基本配置
type RedisConfig struct {
	RedisHost               string        `ini:"redis_host"`
	RedisPassword           string        `ini:"redis_password"`
	RedisDb                 int           `ini:"redis_db"`
	RedisPoolSize           int           `ini:"redis_pool_size"`
	RedisDialTimeout        time.Duration `ini:"redis_dial_timeout"`
	RedisIdleCheckFrequency time.Duration `ini:"redis_idle_check_frequency"`
	RedisIdleTimeout        time.Duration `ini:"redis_idle_timeout"`
	RedisMaxRetries         int           `ini:"redis_max_retries"`
}

//EtcdConfig  etcd服务端的地址列表，用逗号分隔开
type EtcdConfig struct {
	Endpoints string `ini:"endpoints"`
}

//KafkaConfig kafka的基本配置
type KafkaConfig struct {
	Host           string        `ini:"host"`
	Topic          string        `ini:"topic"`
	CommitInterval time.Duration `ini:"commit_interval"`
}

//MachineryConfig  Machinery框架任务调度需要的基本配置
type MachineryConfig struct {
	Broker        string `ini:"broker"`
	DefaultQueue  string `ini:"default_queue"`
	ResultBackend string `ini:"result_backend"`
	Exchange      string `ini:"exchange"`
	ExchangeType  string `ini:"exchange_type"`
	BindingKey    string `ini:"binding_key"`
	TaskQueue     string `ini:"task_queue"`
}

//CacheConfig 缓存配置项
type CacheConfig struct {
	ExpireTime  time.Duration `ini:"expire_time"`
	IsAddRandom bool          `ini:"random"` //是否在过期时间后添加随机数,防止集体失效
	Prefix      string        `ini:"prefix"`
	Enable      bool          `ini:"enable"`
	IsLocal     bool          `int:"local"` //是否启用本地缓存
	MaxNums     int           `ini:"max_nums"`
}

var (
	once           sync.Once
	instance       *GConf
	ErrConfFileNil = errors.New("conf file not init")
)

type GConf struct {
	defaultSectionKey sync.Map //保存默认分区配置
	conf              *ini.File
}

//GetInstance 单例模式获取配置,保证全局唯一
func GetInstance() *GConf {
	once.Do(func() {
		instance = &GConf{}
	})
	return instance
}

//InitConfig 初始化配置
// Deprecated: use toml instead @2019.06.10 by dahe.lai.
func (c *GConf) Init(configFile string) error {
	initconf, err := ini.Load(configFile) //加载配置文件
	if err != nil {
		return err
	}
	c.conf = initconf
	c.conf.BlockMode = false

	//加载默认section key
	c.loadDefaultSection()

	return nil
}

//从默认分区读取配置
func (c *GConf) Get(key string) string {
	return c.GetWithDefault(key, "")
}

func (c *GConf) GetWithDefault(key string, defaultValue string) string {
	value, ok := c.defaultSectionKey.Load(key)
	if ok {
		return strings.Trim(value.(string), " \n\r")
	} else {
		return defaultValue
	}
}

func (c *GConf) loadDefaultSection() {
	keys := c.conf.Section("").Keys()
	for _, key := range keys {
		c.defaultSectionKey.Store(key.Name(), key.Value())
	}
}

//GetDBConfig  获取db的配置 约定大于配置
func (c *GConf) GetDBConfig(db *DbConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section("db_master").MapTo(&db.Master) //解析db master配置
	if err != nil {
		return err
	}
	// 解析db slave配置
	for _, section := range c.conf.Sections() {
		if strings.HasPrefix(section.Name(), "db_slave") {
			dbSlave := BaseDbConfig{}
			err = section.MapTo(&dbSlave)
			if err != nil {
				return err
			}
			db.Slave = append(db.Slave, dbSlave)
		}
	}
	return nil
}

//GetRedisConf  获取redis的配置  约定大于配置
func (c *GConf) GetRedisConf(redis *RedisConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section("redis").MapTo(redis) //解析redis配置
	if err != nil {
		return err
	}
	return nil
}

//GetEtcdConf  获取etcd的配置  约定大于配置
func (c *GConf) GetEtcdConf(etcd *EtcdConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section("etcd").MapTo(etcd) //解析etcd配置
	if err != nil {
		return err
	}
	return nil
}

//GetKafkaConf  获取kafka的配置  约定大于配置
func (c *GConf) GetKafkaConf(kafka *KafkaConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section("kafka").MapTo(kafka) //解析kafka配置
	if err != nil {
		return err
	}
	return nil
}

//GetMachineryConf  获取machinery的配置  约定大于配置
func (c *GConf) GetMachineryConf(machinery *MachineryConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section("machinery").MapTo(machinery) //解析machinery配置
	if err != nil {
		return err
	}
	return nil
}

//GetCacheConf 因为cache存在多个,所以此处因为输入名称
func (c *GConf) GetCacheConf(name string, cache *CacheConfig) error {
	if c.conf == nil {
		return ErrConfFileNil
	}

	err := c.conf.Section(name).MapTo(cache) //解析cache配置
	if err != nil {
		return err
	}
	return nil
}
