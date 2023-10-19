package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cart struct {
	ID              string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	CartItems       []CartItem
	BaseTotalPrice  decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxAmount       decimal.Decimal `gorm:"type:decimal(16,2)"`
	TaxPercent      decimal.Decimal `gorm:"type:decimal(10,2)"`
	DiscountAmount  decimal.Decimal `gorm:"type:decimal(16,2)"`
	DiscountPercent decimal.Decimal `gorm:"type:decimal(10,2)"`
	GrandTotal      decimal.Decimal `gorm:"type:decimal(16,2)"`
	Coupon          int
}

func (c *Cart) GetCart(db *gorm.DB, cartID string) (*Cart, error) {
	var (
		err  error
		cart Cart
	)

	err = db.Debug().Model(Cart{}).Where("id = ?", cartID).First(&cart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (c *Cart) CreateCart(db *gorm.DB, cartID string) (*Cart, error) {

	cart := &Cart{
		ID:              cartID,
		BaseTotalPrice:  decimal.Decimal{},
		TaxAmount:       decimal.Decimal{},
		TaxPercent:      decimal.Decimal{},
		DiscountAmount:  decimal.Decimal{},
		DiscountPercent: decimal.Decimal{},
		GrandTotal:      decimal.Decimal{},
	}

	err := db.Debug().Create(&cart).Error
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *Cart) AddItem(db *gorm.DB, item CartItem) (*CartItem, error) {
	var (
		existItem, updateItem CartItem
		product               Product
	)
	err := db.Debug().Model(Product{}).Where("id = ?", item.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	basePrice, _ := product.Price.Float64()
	taxAmount := GetTaxAmount(basePrice)
	discountAmount := 0.0

	err = db.Debug().Model(CartItem{}).
		Where("cart_id = ?", c.ID).
		Where("product_id = ?", product.ID).
		First(&existItem).Error

	if err != nil {
		subTotal := float64(item.Qty) * (basePrice + taxAmount - discountAmount)

		item.CartID = c.ID
		item.BasePrice = product.Price
		item.BaseTotal = decimal.NewFromFloat(basePrice * float64(item.Qty))
		item.TaxPercent = decimal.NewFromFloat(GetTaxPercent())
		item.TaxAmount = decimal.NewFromFloat(taxAmount)
		item.DiscountPercent = decimal.NewFromFloat(0)
		item.DiscountAmount = decimal.NewFromFloat(discountAmount)
		item.SubTotal = decimal.NewFromFloat(subTotal)

		err = db.Debug().Create(&item).Error
		if err != nil {
			return nil, err
		}

		return &item, nil
	}

	updateItem.Qty = existItem.Qty + item.Qty
	updateItem.BaseTotal = decimal.NewFromFloat(basePrice * float64(updateItem.Qty))

	subTotal := float64(updateItem.Qty) * (basePrice + taxAmount - discountAmount)
	updateItem.SubTotal = decimal.NewFromFloat(subTotal)

	err = db.Debug().First(&existItem, "id = ?", existItem.ID).Updates(updateItem).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (c *Cart) GetItems(db *gorm.DB, cartID string) ([]CartItem, error) {
	var items []CartItem

	err := db.Debug().Preload("Product").Model(&CartItem{}).
		Where("cart_id = ?", cartID).
		Order("created_at desc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
