package errorhandling_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestErrorhandling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errorhandling Suite")
}
