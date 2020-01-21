package configs

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	MysqlConf   *MysqlConfig
	CommentConf *CommentConfig
}

type MysqlConfig struct {
	Addr         string
	Db           string
	User         string
	Password     string
	MaxIdleConns int
	MaxOpenConns int
}

type CommentConfig struct {
	AddCoin int
}

func NewConfig() *Config {
	return &Config{
		MysqlConf: &MysqlConfig{
			Addr:         viper.GetString("mysql.addr"),
			Db:           viper.GetString("mysql.db"),
			User:         viper.GetString("mysql.user"),
			Password:     viper.GetString("mysql.password"),
			MaxIdleConns: viper.GetInt("mysql.maxIdleConns"),
			MaxOpenConns: viper.GetInt("mysql.maxOpenConns"),
		},
		CommentConf: &CommentConfig{
			AddCoin: viper.GetInt("comment.addcoin"),
		},
	}
}

func init() {
	if err := initConfig(); err != nil {
		panic(err)
	}
}

func initConfig() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("post-srv")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
