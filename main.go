package main

import (
	IRedis "go4/internal/redis"
	redis "go4/internal/redis/driver"

	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Dir      string `env:"CONFIG_DIR" envDefault:"config/config.json"`
	Port     string
	LogType  string
	LogLevel string
	LogFile  string
	LogAddr  string
	DB       string
	DBConfig
}

type DBConfig struct {
	Driver          string
	Host            string
	Username        string
	Password        string
	Database		string
	SslMode         string
	Timeout         string
	MaxOpenConns    string
	MaxIdeConns     string
	ConnMaxLifetime string
}

var config Config

func init() {
	// Parse cac gia tri cua config xem dung chua
	if err := env.Parse(&config); err != nil {
		log.Error("Error get config value fail")
		log.Fatal(err)
	}
	// Xac dinh path file config => config/config.json
	viper.SetConfigFile(config.Dir)
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err.Error())
		panic(err)
	}

	configure := Config{
		Dir: config.Dir,
		Port: viper.GetString(`main.port`),
		LogType: viper.GetString(`main.log_type`),
		LogLevel: viper.GetString(`main.log_level`),
		LogFile: viper.GetString(`main.log_file`),
		LogAddr: viper.GetString(`main.log_addr`),
		DB: viper.GetString(`main.db`),
	}

	if configure.DB == "enabled" {
		configure.DBConfig = DBConfig{
			Driver: viper.GetString(`db.driver`),
			Host: viper.GetString(`db.host`),
			Username: viper.GetString(`db.username`),
			Password: viper.GetString(`db.password`),
			Database: viper.GetString(`db.database`),
			SslMode: viper.GetString(`db.ssl_mode`),
			Timeout: viper.GetString(`db.timeout`),
			MaxOpenConns: viper.GetString(`db.max_open_conns`),
			MaxIdeConns: viper.GetString(`db.max_ide_conns`),
			ConnMaxLifetime: viper.GetString(`db.conn_max_lifetime`),
		}
	}

	var err error
	IRedis.Redis, err := redis.NewRedis(redis.Config{
		Addr: viper.GetString(`redis.addr`),
		DB: viper.GetString(`redis.database`),
		Password: viper.GetString(`redis.password`),
	})

	
}