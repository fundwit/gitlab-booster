package persistence_test

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab-booster/config"
	"gitlab-booster/persistence"
)

var _ = Describe("ParseDatabaseConfigFromEnv", func() {
	BeforeEach(func() {
		// 重置环境变量值
		os.Setenv(config.EnvDatabaseURL, "")
	})

	It("should return a DatabaseConfig instance when environment variable is valid", func() {
		os.Setenv(config.EnvDatabaseURL, "Mysql://user:pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true")

		config, err := persistence.ParseDatabaseConfigFromEnv()

		Expect(err).To(BeNil())
		Expect(config).NotTo(BeNil())
		Expect(config.DriverType).To(Equal("mysql"))
		Expect(config.DriverArgs).To(Equal("user:pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true"))
	})

	It("should return error when environment variable is not valid", func() {
		os.Setenv(config.EnvDatabaseURL, "user:pwd@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=true")

		config, err := persistence.ParseDatabaseConfigFromEnv()

		Expect(err).NotTo(BeNil())
		Expect(config).To(BeNil())
		Expect(strings.Contains(err.Error(), "environment value is not valid")).To(BeTrue())
	})
})
