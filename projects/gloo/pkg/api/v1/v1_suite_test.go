package v1

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAttributeRoleUpstreamVirtualService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AttributeRoleUpstreamVirtualService Suite")
}
