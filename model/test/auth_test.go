package test

import (
	"payhere/model"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Auth")
}

var _ = Describe("Model Auth Test", func() {
	var (
		auth *model.UserAuth
	)
	BeforeEach(func() {

	})
	It("would be valid", func() {
		auth = &model.UserAuth{
			Phone:    "010-0000-0000",
			Password: "a123",
		}
		Expect(auth.Validate()).To(Equal(true))
	})
	It("would be invalid", func() {
		auth = &model.UserAuth{
			Phone:    "010-abc-0000",
			Password: "a123",
		}
		Expect(auth.Validate()).To(Equal(false))
	})
})
