package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cart struct {
	ID              string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	Items           []Item
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

	err = db.Debug().Preload("Items").Model(Cart{}).Where("id = ?", cartID).First(&cart).Error
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
		Coupon:          0,
	}

	err := db.Debug().Create(&cart).Error
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *Cart) AddItem(db *gorm.DB, item Item) (*Item, error) {
	var (
		existItem, updateItem Item
		product               Product
	)

	err := db.Debug().Model(Product{}).Where("id = ?", item.ProductID).First(&product).Error
	if err != nil {
		return nil, err
	}

	basePrice, _ := product.Price.Float64()
	taxAmount := GetTaxAmount(basePrice)
	discount := 0.0
	// coupon := 0

	err = db.Debug().Model(Item{}).
		Where("cart_id = ?", c.ID).
		Where("product_id = ?", product.ID).
		First(&existItem).Error
	if err != nil {
		subTotal := float64(item.Qty) * (basePrice + taxAmount - discount)

		item.CartID = c.ID
		item.BasePrice = product.Price
		item.BaseTotal = decimal.NewFromFloat(basePrice * float64(item.Qty))
		item.TaxPercent = decimal.NewFromFloat(GetTaxPercent())
		item.TaxAmount = decimal.NewFromFloat(taxAmount)
		item.DiscountPercent = decimal.NewFromFloat(discount)

		item.SubTotal = decimal.NewFromFloat(subTotal)

		err = db.Debug().Create(&item).Error
		if err != nil {
			return nil, err
		}
		return &item, nil
	}

	updateItem.Qty = existItem.Qty + item.Qty
	updateItem.BaseTotal = decimal.NewFromFloat(basePrice * float64(updateItem.Qty))

	subTotal := float64(updateItem.Qty) * (basePrice + taxAmount - discount)
	updateItem.SubTotal = decimal.NewFromFloat(subTotal)

	err = db.Debug().Model(&Item{}).Where("id = ?", existItem.ID).Updates(updateItem).Error
	if err != nil {
		return nil, err
	}

	return &item, nil

}

func (c *Cart) CalculateCart(db *gorm.DB, cartID string) (*Cart, error) {
	cartBaseTotalPrice := 0.0
	cartTaxAmount := 0.0
	cartDiscountAmount := 0.0
	cartGrandTotal := 0.0
	// Coupon := 0

	for _, item := range c.Items {
		itemBaseTotal, _ := item.BaseTotal.Float64()
		itemTaxAmount, _ := item.TaxAmount.Float64()
		itemSubTotalTaxAmount := itemTaxAmount * float64(item.Qty)
		itemDiscountAmount, _ := item.DiscountAmount.Float64()
		itemSubTotalDiscountAmount := itemDiscountAmount * float64(item.Qty)
		itemSubTotal, _ := item.SubTotal.Float64()

		cartBaseTotalPrice += itemBaseTotal
		cartTaxAmount += itemSubTotalTaxAmount
		cartDiscountAmount += itemSubTotalDiscountAmount
		cartGrandTotal += itemSubTotal

	}

	var updateCart, cart Cart

	if cart.GrandTotal.GreaterThan(decimal.NewFromInt(50000)) {
		updateCart.Coupon++
	}

	if cart.GrandTotal.GreaterThan(decimal.NewFromInt(100000)) {
		// Calculate how many multiples of 100,000 fit in updateItem.BaseTotal
		multiplesCoupon := int(updateCart.GrandTotal.Div(decimal.NewFromInt(100000)).IntPart())
		updateCart.Coupon += multiplesCoupon
	}

	updateCart.BaseTotalPrice = decimal.NewFromFloat(cartBaseTotalPrice)
	updateCart.TaxAmount = decimal.NewFromFloat(cartTaxAmount)
	updateCart.DiscountAmount = decimal.NewFromFloat(cartDiscountAmount)
	updateCart.GrandTotal = decimal.NewFromFloat(cartGrandTotal)
	// updateCart.Coupon = Coupon

	err := db.Debug().Model(&Cart{}).Where("id = ?", c.ID).Updates(updateCart).Error
	if err != nil {
		return nil, err
	}

	return &cart, nil
}
