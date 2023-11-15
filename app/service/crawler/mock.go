package crawler

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

type Mock struct {
	size int
}

func NewMock(size int) *Mock {
	return &Mock{size}
}

var _ Crawler = (*Mock)(nil)

func (m *Mock) SellerType() SellerType {
	return SellerTypeOther
}

func (m *Mock) Fetch(ctx context.Context) ([]Product, error) {
	fakeProducts := make([]Product, 0, m.size)

	for i := 1; i <= m.size; i++ {
		product := Product{
			ID:          uint64(i),
			Name:        fmt.Sprintf("Product %d", i),
			Description: fmt.Sprintf("Description for Product %d", i),
			Type:        ProductType(i % (int(ProductTypeEarphone) + 1)),
			UpdatedAt:   time.Now(),
			ProductPrice: ProductPrice{
				Price:      genRandDecimal(1, 1000),
				Discount:   genRandDecimal(0, 100),
				SourceURL:  fmt.Sprintf("https://example.com/product/%d", i),
				SellerType: SellerType(i % (int(SellerTypeOther) + 1)),
			},
		}

		fakeProducts = append(fakeProducts, product)
	}

	return fakeProducts, nil
}

func genRandDecimal(min, max int64) decimal.Decimal {
	// Generate a random int64 within the specified range.
	randInt, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	randomInt64 := randInt.Int64() + min
	return decimal.NewFromInt(randomInt64)
}
