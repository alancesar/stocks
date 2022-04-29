package repositoty

import (
	"context"
	"gorm.io/gorm"
	"stocks/currency"
	"stocks/operation"
	"stocks/stock"
)

type (
	GormDatabase struct {
		DB *gorm.DB
	}

	entity struct {
		gorm.Model
		operation.Operation
	}

	entry struct {
		Stock        stock.Stock
		Quantity     int
		AveragePrice float64
		CurrentPrice float64
		Investment   float64
		Settled      float64
	}

	entries []entry
)

func (e entity) TableName() string {
	return "operations"
}

func (s entries) ToDomain() operation.Summary {
	var summary operation.Summary

	for _, e := range s {
		summary = append(summary, operation.Entry{
			Stock:        e.Stock,
			Quantity:     e.Quantity,
			AveragePrice: currency.NewFromFloat(e.AveragePrice),
			LastPrice:    currency.NewFromFloat(e.CurrentPrice),
			Investment:   currency.NewFromFloat(e.Investment),
			Settled:      currency.NewFromFloat(e.Settled),
		})
	}

	return summary
}

func NewGormDatabase(db *gorm.DB) *GormDatabase {
	_ = db.AutoMigrate(&entity{})
	return &GormDatabase{
		DB: db,
	}
}

func (d GormDatabase) Create(ctx context.Context, op operation.Operation) error {
	query := d.DB.WithContext(ctx).Create(&entity{
		Operation: op,
	})
	return query.Error
}
func (d GormDatabase) Summary(ctx context.Context) (operation.Summary, error) {
	var e entries

	query := d.DB.WithContext(ctx).Raw(`
		SELECT buy.stock,
			   buy.total_quantity - IFNULL(sell.total_quantity, 0) quantity,
			   buy.average_price                                   average_price,
			   buy.total_amount                                    investment,
			   sell.total_amount                                   settled
		FROM (SELECT stock                                                stock,
					 round(sum(quantity * unit_value), 2)                 total_amount,
					 sum(quantity)                                        total_quantity,
					 round(sum(quantity * unit_value) / sum(quantity), 2) average_price
			  FROM operations
			  WHERE type = ?
			  GROUP BY stock
			  ORDER BY stock) buy
				 LEFT JOIN (SELECT stock                                stock,
								   round(sum(quantity * unit_value), 2) total_amount,
								   sum(quantity)                        total_quantity
							FROM operations
							WHERE type = ?
							GROUP BY stock) sell
						   ON sell.stock == buy.stock;
	`, operation.Buy, operation.Sell).Scan(&e)

	if query.Error != nil {
		return nil, query.Error
	}

	return e.ToDomain(), nil
}
