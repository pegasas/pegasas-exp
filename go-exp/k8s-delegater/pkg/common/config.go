package common

import (
	"github.com/spf13/viper"
	rawlog "log"
)

var config *viper.Viper

const (
	SQL_SERVER_HOST     = "sql.server.host"
	SQL_SERVER_PORT     = "sql.server.port"
	SQL_SERVER_DATABASE = "sql.server.database"
	SQL_SERVER_USERNAME = "sql.server.username"
)

func InitConfig() {
	var err error
	env := GetEnv()
	rawlog.Print(env)
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func GetConfig() *viper.Viper {
	return config
}

func GetSQLServerHost() string {
	config := GetConfig()
	return config.GetString(SQL_SERVER_HOST)
}

func GetSQLServerPort() int {
	config := GetConfig()
	return config.GetInt(SQL_SERVER_PORT)
}

func GetSQLServerDatabase() string {
	config := GetConfig()
	return config.GetString(SQL_SERVER_DATABASE)
}

func GetSQLServerUsername() string {
	config := GetConfig()
	return config.GetString(SQL_SERVER_USERNAME)
}
