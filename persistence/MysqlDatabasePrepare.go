package persistence

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// PrepareMysqlDatabase 自动创建 MySql 数据库， 输出参数格式为：'user:pwd@tcp(host:3306)/dbname?xxx=xxx'
func PrepareMysqlDatabase(mysqlDriverArgs string) error {
	databaseName, rootDriverArgs, err := extractDatabaseName(mysqlDriverArgs)
	if err != nil {
		return err
	}

	db, err := gorm.Open("mysql", rootDriverArgs)
	if err != nil {
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		return err
	}

	initSQL := "CREATE DATABASE IF NOT EXISTS `" + databaseName + "` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;"
	err = db.Exec(initSQL).Error
	if err != nil {
		return err
	}

	return nil
}

func extractDatabaseName(mysqlDriverArgs string) (string, string, error) {
	nameIndex := strings.IndexRune(mysqlDriverArgs, '/')
	paramsIndex := strings.IndexRune(mysqlDriverArgs, '?')

	if nameIndex > 0 && paramsIndex > nameIndex {
		return mysqlDriverArgs[nameIndex+1 : paramsIndex], mysqlDriverArgs[0:nameIndex+1] + mysqlDriverArgs[paramsIndex:], nil
	}
	if nameIndex < 0 {
		return "", mysqlDriverArgs, nil
	}
	if nameIndex > 0 && paramsIndex < 0 {
		return mysqlDriverArgs[nameIndex+1:], mysqlDriverArgs[0 : nameIndex+1], nil
	}
	return "", "", errors.New("invalid mysql driver args")
}
