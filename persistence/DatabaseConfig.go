package persistence

import (
	"errors"
	"gitlab-booster/config"
	"os"
	"strings"
)

// DatabaseConfig 数据库连接配置信息
type DatabaseConfig struct {
	DriverType string
	DriverArgs string
}

// ParseDatabaseConfigFromEnv 从环境变量获取数据库连接配置信息
func ParseDatabaseConfigFromEnv() (*DatabaseConfig, error) {
	// mysql://${MYSQL_USER:archivist}:${MYSQL_PWD:not_specified}@tcp(${MYSQL_SERVICE:os-mysql-svc.default:3306})/${MYSQL_DBNAME:cicd_archivist}?charset=utf8mb4&parseTime=true&loc=Local
	databaseURL := os.ExpandEnv(os.Getenv(config.EnvDatabaseURL))
	slice := strings.Split(databaseURL, "://")
	if len(slice) != 2 {
		return nil, errors.New(config.EnvDatabaseURL + " environment value is not valid, a correct example: " +
			"'mysql://user:pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true'")
	}

	return &DatabaseConfig{DriverType: strings.ToLower(slice[0]), DriverArgs: slice[1]}, nil
}
