package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"stocks/asset"
	"stocks/currency"
	"stocks/operation"
	"stocks/stock"
)

type (
	GormDatabase struct {
		DB *gorm.DB
	}

	operationEntity struct {
		gorm.Model
		operation.Operation
	}

	detailsEntity struct {
		gorm.Model
		stock.Details
	}

	assetEntity struct {
		Symbol       stock.Symbol
		Quantity     int
		AveragePrice float64
		CurrentPrice float64
		Investment   float64
		Settled      float64
	}

	assets []assetEntity
)

func (e operationEntity) TableName() string {
	return "operations"
}

func (e detailsEntity) TableName() string {
	return "details"
}

func (s assets) ToDomain() asset.Assets {
	var a asset.Assets

	for _, e := range s {
		a = append(a, asset.Asset{
			Symbol:       e.Symbol,
			Quantity:     e.Quantity,
			AveragePrice: currency.NewFromFloat(e.AveragePrice),
			LastPrice:    currency.NewFromFloat(e.CurrentPrice),
			Investment:   currency.NewFromFloat(e.Investment),
			Settled:      currency.NewFromFloat(e.Settled),
		})
	}

	return a
}

func NewGormDatabase(db *gorm.DB) *GormDatabase {
	_ = db.AutoMigrate(&operationEntity{})
	return &GormDatabase{
		DB: db,
	}
}

func (d GormDatabase) Create(ctx context.Context, op operation.Operation) error {
	return d.DB.WithContext(ctx).Create(&operationEntity{
		Operation: op,
	}).Error
}

func (d GormDatabase) List(ctx context.Context) (operation.List, error) {
	var entities []operationEntity
	query := d.DB.WithContext(ctx).Raw("SELECT * FROM operations ORDER BY date, id").Scan(&entities)
	if query.Error != nil {
		return nil, query.Error
	}

	operations := make(operation.List, len(entities))
	for i := range entities {
		operations[i] = entities[i].Operation
	}

	return operations, nil
}

func (d GormDatabase) Assets(ctx context.Context) (asset.Assets, error) {
	var a assets

	query := d.DB.WithContext(ctx).Raw(`
		SELECT buy.symbol                                          symbol,
			   buy.total_quantity - IFNULL(sell.total_quantity, 0) quantity,
			   buy.average_price                                   average_price,
			   buy.total_amount                                    investment,
			   sell.total_amount                                   settled
		FROM (SELECT symbol                                               symbol,
					 round(sum(quantity * unit_value), 2)                 total_amount,
					 sum(quantity)                                        total_quantity,
					 round(sum(quantity * unit_value) / sum(quantity), 2) average_price
			  FROM operations
			  WHERE type = ?
			  GROUP BY symbol
			  ORDER BY symbol) buy
				 LEFT JOIN (SELECT symbol                               symbol,
								   round(sum(quantity * unit_value), 2) total_amount,
								   sum(quantity)                        total_quantity
							FROM operations
							WHERE type = ?
							GROUP BY symbol) sell
						   ON sell.symbol == buy.symbol;
	`, operation.Buy, operation.Sell).Scan(&a)

	if query.Error != nil {
		return nil, query.Error
	}

	return a.ToDomain(), nil
}

func (d GormDatabase) GetDetails(ctx context.Context, symbol stock.Symbol) (stock.Details, error) {
	var entity detailsEntity
	if query := d.DB.WithContext(ctx).Find(&entity, "symbol = ?", symbol); query.Error != nil {
		return stock.Details{}, query.Error
	} else if query.RowsAffected == 0 {
		return stock.Details{}, errors.New("not found")
	}

	return stock.Details{
		Symbol:  entity.Symbol,
		Type:    entity.Type,
		Sector:  entity.Sector,
		Segment: entity.Segment,
	}, nil
}

func (d GormDatabase) InsertDetails(ctx context.Context, details stock.Details) error {
	return d.DB.WithContext(ctx).Create(&detailsEntity{
		Details: details,
	}).Error
}
