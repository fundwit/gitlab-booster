package testinfra

import (
	"gitlab-booster/persistence"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TestDatabase struct {
	TestDatabaseName string
	DS               *persistence.DatasourceManager
}

// StartMysqlTestDatabase TEST_MYSQL_SERVICE=root:root@(127.0.0.1:3306)
func StartMysqlTestDatabase(baseName string) *TestDatabase {
	mysqlSvc := os.Getenv("TEST_MYSQL_SERVICE")
	if mysqlSvc == "" {
		mysqlSvc = "root:root@(127.0.0.1:3306)"
	}
	databaseName := baseName + "_test_" + strings.ReplaceAll(uuid.New().String(), "-", "")

	dbConfig := &persistence.DatabaseConfig{
		DriverType: "mysql", DriverArgs: mysqlSvc + "/" + databaseName + "?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s",
	}

	// create database (no conflict)
	if err := persistence.PrepareMysqlDatabase(dbConfig.DriverArgs); err != nil {
		log.Fatalf("failed to prepare database %v\n", err)
	}

	ds := &persistence.DatasourceManager{DatabaseConfig: dbConfig}
	// connect
	if err := ds.Start(); err != nil {
		defer ds.Stop()
		log.Fatalf("database conneciton failed %v\n", err)
	}

	return &TestDatabase{TestDatabaseName: databaseName, DS: ds}
}

func StopMysqlTestDatabase(testDatabase *TestDatabase) {
	if testDatabase != nil || testDatabase.DS != nil {
		if testDatabase.DS.GromDB() != nil {
			if err := testDatabase.DS.GromDB().Exec("DROP DATABASE " + testDatabase.TestDatabaseName).Error; err != nil {
				log.Println("failed to drop test database: " + testDatabase.TestDatabaseName)
			} else {
				log.Println("test database " + testDatabase.TestDatabaseName + " dropped")
			}
		}

		// close connection
		testDatabase.DS.Stop()
	}
}
