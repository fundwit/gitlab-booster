package main

import (
	"fmt"
	"gitlab-booster/config"
	"gitlab-booster/persistence"
	"log"
	"os"
	"strings"
)

func main() {
	config.CheckConfig()

	initMode := strings.ToLower(os.Getenv("INIT_MODE")) == "true"
	fmt.Printf("gitlab-booster service, init mode: %v\n", initMode)

	dbConfig, err := persistence.ParseDatabaseConfigFromEnv()
	if err != nil {
		log.Fatalf("parse database config failed %v\n", err)
	}

	if err = InitializeService(dbConfig); err != nil {
		log.Fatal(err)
	}

	// 如果仅初始化，则退出执行
	if initMode {
		return
	}

	ServerHTTP(dbConfig)
}
