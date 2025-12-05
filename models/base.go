package models

import (
	"context"
	"time"

	"zatrano/pkg/currentuser"

	"gorm.io/gorm"
)

// BaseModel tüm modeller için ortak yapı
type BaseModel struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreatedBy uint  `gorm:"column:created_by;index"`
	UpdatedBy uint  `gorm:"column:updated_by;index"`
	DeletedBy *uint `gorm:"column:deleted_by;index"`
	IsActive  bool  `gorm:"default:true;index"`
}

// helper: context'ten user ID al
func getCurrentUserID(ctx context.Context) uint {
	cu := currentuser.FromContext(ctx)
	if cu.ID != 0 {
		return cu.ID
	}
	return 0
}

// RegisterBaseModelCallbacks tüm DB için global callback olarak eklenir
func RegisterBaseModelCallbacks(db *gorm.DB) {

	// Create işlemleri
	db.Callback().Create().Before("gorm:create").Register("base_model:before_create", func(tx *gorm.DB) {
		cuID := getCurrentUserID(tx.Statement.Context)
		if cuID == 0 || tx.Statement.Schema == nil {
			return
		}

		rv := tx.Statement.ReflectValue

		if f := tx.Statement.Schema.LookUpField("created_by"); f != nil {
			_ = f.Set(tx.Statement.Context, rv, cuID)
		}
		if f := tx.Statement.Schema.LookUpField("updated_by"); f != nil {
			_ = f.Set(tx.Statement.Context, rv, cuID)
		}
		if f := tx.Statement.Schema.LookUpField("is_active"); f != nil {
			if !rv.FieldByName("IsActive").Bool() {
				_ = f.Set(tx.Statement.Context, rv, true)
			}
		}
	})

	// Update işlemleri
	db.Callback().Update().Before("gorm:update").Register("base_model:before_update", func(tx *gorm.DB) {
		cuID := getCurrentUserID(tx.Statement.Context)
		if cuID == 0 || tx.Statement.Schema == nil {
			return
		}

		rv := tx.Statement.ReflectValue
		if f := tx.Statement.Schema.LookUpField("updated_by"); f != nil {
			_ = f.Set(tx.Statement.Context, rv, cuID)
		}
	})

	// Soft Delete işlemleri
	db.Callback().Delete().Before("gorm:delete").Register("base_model:before_delete", func(tx *gorm.DB) {
		cuID := getCurrentUserID(tx.Statement.Context)
		if cuID == 0 || tx.Statement.Schema == nil {
			return
		}

		rv := tx.Statement.ReflectValue
		if f := tx.Statement.Schema.LookUpField("deleted_by"); f != nil {
			_ = f.Set(tx.Statement.Context, rv, cuID)
		}
		if f := tx.Statement.Schema.LookUpField("updated_by"); f != nil {
			_ = f.Set(tx.Statement.Context, rv, cuID)
		}
	})
}
