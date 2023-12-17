package test

import (
	"context"
	"math/rand"
	"payhere/model"
	"payhere/util"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestProduct(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Product")
}

var _ = Describe("Model Product Test", func() {
	var (
		product *model.Product
		ctx     context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
	})
	It("would be create invalid", func() {
		product = &model.Product{
			UID:      uuid.NewString(),
			ShopID:   rand.Uint64(),
			State:    model.ProductStateNew,
			Category: "category",
			Price:    rand.Uint64(),
		}
		Expect(product.CreateValidate(ctx)).To(Equal(false))
	})
	It("would be create valid", func() {
		originalPrice := rand.Uint64()
		barcode := "barcode"
		description := "description"
		expireDate := time.Now().UTC()
		uid := uuid.NewString()
		ctx = context.WithValue(ctx, util.OwnerKey, uid)
		product = &model.Product{
			UID:           uid,
			Name:          "name",
			ShopID:        rand.Uint64(),
			State:         model.ProductStateNew,
			Category:      "category",
			Price:         rand.Uint64(),
			OriginalPrice: &originalPrice,
			Description:   &description,
			Barcode:       &barcode,
			ExpireDate:    &expireDate,
			Size:          model.ProductSizeLarge,
		}
		Expect(product.CreateValidate(ctx)).To(Equal(true))
	})
	It("would be update invalid", func() {
		product = &model.Product{
			UID: uuid.NewString(),
		}
		Expect(product.UpdateValidate(ctx)).To(Equal(false))
	})
	It("would be update valid", func() {
		uid := uuid.NewString()
		ctx = context.WithValue(ctx, util.OwnerKey, uid)
		product = &model.Product{
			UID:  uid,
			Name: "change name",
		}
		Expect(product.UpdateValidate(ctx)).To(Equal(true))
	})
})
