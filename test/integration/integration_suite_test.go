package integration_test

import (
	// "k8s.io/kubernetes/test/e2e/framework"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	// SynchronizedBeforeSuite(nil, func(data []byte) {})
}

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}
