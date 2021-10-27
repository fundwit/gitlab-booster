package main

import (
	"fmt"
	"gitlab-booster/errorhandling"
	"gitlab-booster/persistence"
	"gitlab-booster/servehttp"
	"log"

	"github.com/gin-gonic/gin"
)

// InitializeService 服务初始化
func InitializeService(dbConfig *persistence.DatabaseConfig) error {
	// 创建数据库（多实例执行不会冲突）
	if dbConfig.DriverType == "mysql" {
		log.Printf("database type is %s, create database if not exist\n", dbConfig.DriverType)
		if err := persistence.PrepareMysqlDatabase(dbConfig.DriverArgs); err != nil {
			return fmt.Errorf("failed to prepare mysql database: %v", err)
		}
	}

	// 连接数据库
	ds := &persistence.DatasourceManager{DatabaseConfig: dbConfig}
	if err := ds.Start(); err != nil {
		return fmt.Errorf("database conneciton failed: %v", err)
	}
	defer ds.Stop()

	// 数据库升级 (TODO 分布式锁限制单实例执行)
	err := ds.GromDB().AutoMigrate(
		&servehttp.Manifest{}, &servehttp.RepositoryRef{},
	).Error
	if err != nil {
		return fmt.Errorf("database migration failed: %v", err)
	}

	return nil
}

// ServerHTTP 启动服务
func ServerHTTP(dbConfig *persistence.DatabaseConfig) {
	// 连接数据库
	ds := &persistence.DatasourceManager{DatabaseConfig: dbConfig}
	if err := ds.Start(); err != nil {
		log.Fatalf("database conneciton failed %v\n", err)
	}
	defer ds.Stop()

	//archiveSvc := archive.NewArchiveMetaService(ds, 0)
	//securityContextResolver := security.NewSecurityContextResolver()

	// 启动 web 服务
	router := gin.New()
	router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/"}}),
		gin.Recovery(),
		errorhandling.GolbalErrorHandlingFilter(),
	)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "OK"})
	})

	servehttp.RegisterGitlabV4AuthHandlers(router)
	servehttp.RegisterManifestHandlers(router, ds)

	//servehttp.RegisterMetadataHandlers(router)
	//servehttp.RegisterArchiveHandlers(router, archiveSvc, servehttp.NewSecurityContextMustValidFilter(securityContextResolver))

	servehttp.StartHTTPServer(router)
}
