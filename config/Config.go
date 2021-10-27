package config

import (
	"errors"
	"os"
)

// CheckConfig 验证配置
func CheckConfig() {
	GetGitlabEndpoint()
	GetGitlabOAuthClientID()
	GetGitlabOAuthClientSecret()
}

// GetGitlabEndpoint  读取配置：请求 gitlab 服务的地址
func GetGitlabEndpoint() string {
	value := os.Getenv("GITLAB_ENDPOINT")
	if value == "" {
		panic(errors.New("GITLAB_ENDPOINT env is empty"))
	}
	return value
}

// GetGitlabOAuthClientID  读取配置: OAuth client id
func GetGitlabOAuthClientID() string {
	value := os.Getenv("GITLAB_OAUTH_CLIENT_ID")
	if value == "" {
		panic(errors.New("GITLAB_OAUTH_CLIENT_ID env is empty"))
	}
	return value
}

// GetGitlabOAuthClientSecret  读取配置: OAuth client secret
func GetGitlabOAuthClientSecret() string {
	value := os.Getenv("GITLAB_OAUTH_CLIENT_SECRET")
	if value == "" {
		panic(errors.New("GITLAB_OAUTH_CLIENT_SECRET env is empty"))
	}
	return value
}

// EnvDatabaseURL 保存数据库连接配置的环境变量名称。这个变量中可以包含其他环境变量
// 示例 mysql://${MYSQL_USER:archivist}:${MYSQL_PWD:not_specified}@tcp(${MYSQL_SERVICE:os-mysql-svc.default:3306})/${MYSQL_DBNAME:cicd_archivist}?charset=utf8mb4&parseTime=true&loc=Local
const EnvDatabaseURL = "DATABASE_URL"
