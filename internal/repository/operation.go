package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"stocks/asset"
	"stocks/currency"
	"stocks/operation"
	"stocks/stock"
	"time"
)

type (
	GormDatabase struct {
		DB *gorm.DB
	}

	Operation struct {
		gorm.Model
		Symbol    string
		Type      int
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	Detail struct {
		gorm.Model
		Symbol    string
		Name      string
		Sector    string
		SubSector string
		Segment   string
	}

	Asset struct {
		Symbol       stock.Symbol
		Quantity     int
		AveragePrice float64
		CurrentPrice float64
		Investment   float64
		Settled      float64
	}

	Assets []Asset
)

func (s Assets) ToDomain() asset.Assets {
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
	_ = db.AutoMigrate(&Operation{}, &Detail{})

	return &GormDatabase{
		DB: db,
	}
}

func (d GormDatabase) Create(ctx context.Context, op operation.Operation) error {
	return d.DB.WithContext(ctx).Create(&Operation{
		Symbol:    string(op.Symbol),
		Type:      int(op.Type),
		Quantity:  op.Quantity,
		UnitValue: op.UnitValue,
		Date:      op.Date,
	}).Error
}

func (d GormDatabase) List(ctx context.Context) (operation.List, error) {
	var entities []Operation
	query := d.DB.WithContext(ctx).Raw("SELECT * FROM operations ORDER BY date, id").Scan(&entities)
	if query.Error != nil {
		return nil, query.Error
	}

	operations := make(operation.List, len(entities))
	for i, e := range entities {
		operations[i] = operation.Operation{
			Symbol:    stock.Symbol(e.Symbol),
			Type:      operation.Type(e.Type),
			Quantity:  e.Quantity,
			UnitValue: e.UnitValue,
			Date:      e.Date,
		}
	}

	return operations, nil
}

func (d GormDatabase) Assets(ctx context.Context) (asset.Assets, error) {
	var a Assets

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
	var entity Detail
	if query := d.DB.WithContext(ctx).Find(&entity, "symbol = ?", symbol); query.Error != nil {
		return stock.Details{}, query.Error
	} else if query.RowsAffected == 0 {
		return stock.Details{}, errors.New("not found")
	}

	return stock.Details{
		Symbol:    stock.Symbol(entity.Symbol),
		Name:      entity.Name,
		Sector:    entity.Sector,
		SubSector: entity.SubSector,
		Segment:   entity.Segment,
	}, nil
}

func (d GormDatabase) InsertDetails(ctx context.Context, details stock.Details) error {
	return d.DB.WithContext(ctx).Create(&Detail{
		Symbol:    string(details.Symbol),
		Name:      details.Name,
		Sector:    details.Sector,
		SubSector: details.SubSector,
		Segment:   details.Segment,
	}).Error
}
