package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCopyenv(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Copyenv Suite")
}
