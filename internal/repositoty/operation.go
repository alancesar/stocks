package repositoty

import (
	"context"
	"gorm.io/gorm"
	"stocks/operation"
)

type (
	GormDatabase struct {
		DB *gorm.DB
	}

	entity struct {
		gorm.Model
		operation.Operation
	}
)

func (e entity) TableName() string {
	return "operations"
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
