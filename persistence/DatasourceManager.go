package persistence

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DatasourceManager 管理 grom Datasource
type DatasourceManager struct {
	gormDB         *gorm.DB
	DatabaseConfig *DatabaseConfig
}

// Start start grom datasoruce
func (m *DatasourceManager) Start() error {
	db, err := connect(m.DatabaseConfig.DriverType, m.DatabaseConfig.DriverArgs)
	if err != nil {
		return err
	}

	m.gormDB = db
	if os.Getenv("GIN_MODE") != "release" {
		m.gormDB.LogMode(true)
	}

	return nil
}

// Stop stop gorm datasoruce
func (m *DatasourceManager) Stop() {
	if m.gormDB != nil {
		if err := m.gormDB.Close(); err != nil {
			log.Printf("failed to close DB: %v", err)
		}
		m.gormDB = nil
	}
}

// GromDB 获取 DatasourceManager 中的 grom DB 对象
func (m *DatasourceManager) GromDB() *gorm.DB {
	if m.gormDB != nil {
		return m.gormDB.New()
	}
	return nil
}

func connect(driver, driverArgs string) (*gorm.DB, error) {
	db, err := gorm.Open(driver, driverArgs)
	if err != nil {
		return nil, err
	}
	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
