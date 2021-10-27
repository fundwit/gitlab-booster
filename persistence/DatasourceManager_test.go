package persistence_test

import (
	"gitlab-booster/persistence"
	"gitlab-booster/testinfra"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DatasourceManager", func() {
	Describe("Start and Initilaize Database", func() {
		var testDatabase *testinfra.TestDatabase

		BeforeEach(func() {
			testDatabase = testinfra.StartMysqlTestDatabase("code_superviser")
		})
		AfterEach(func() {
			testinfra.StopMysqlTestDatabase(testDatabase)
		})

		Context("mysql database", func() {
			It("should connect to database when start", func() {
				ds := &persistence.DatasourceManager{
					DatabaseConfig: testDatabase.DS.DatabaseConfig,
				}
				Expect(ds.GromDB()).To(BeNil())

				if err := ds.Start(); err != nil {
					log.Fatal(err)
				}
				Expect(ds.GromDB()).ToNot(BeNil())

				ds.Stop()
				Expect(ds.GromDB()).To(BeNil())
			})
		})
	})
})
